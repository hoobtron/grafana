package authnimpl

import (
	"context"

	"github.com/grafana/grafana/pkg/infra/log"
	"github.com/grafana/grafana/pkg/infra/remotecache"
	"github.com/grafana/grafana/pkg/login/social"
	"github.com/grafana/grafana/pkg/plugins/manager/registry"
	"github.com/grafana/grafana/pkg/services/accesscontrol"
	"github.com/grafana/grafana/pkg/services/apikey"
	"github.com/grafana/grafana/pkg/services/auth"
	"github.com/grafana/grafana/pkg/services/auth/gcomsso"
	"github.com/grafana/grafana/pkg/services/authn"
	"github.com/grafana/grafana/pkg/services/authn/authnimpl/sync"
	"github.com/grafana/grafana/pkg/services/authn/clients"
	"github.com/grafana/grafana/pkg/services/featuremgmt"
	"github.com/grafana/grafana/pkg/services/ldap/service"
	"github.com/grafana/grafana/pkg/services/login"
	"github.com/grafana/grafana/pkg/services/loginattempt"
	"github.com/grafana/grafana/pkg/services/oauthtoken"
	"github.com/grafana/grafana/pkg/services/org"
	"github.com/grafana/grafana/pkg/services/pluginsintegration/pluginsettings"
	"github.com/grafana/grafana/pkg/services/quota"
	"github.com/grafana/grafana/pkg/services/rendering"
	"github.com/grafana/grafana/pkg/services/user"
	"github.com/grafana/grafana/pkg/setting"
)

const cloudHomePluginId = "cloud-home-app"

type Registration struct{}

func ProvideRegistration(
	cfg *setting.Cfg, authnSvc authn.Service,
	orgService org.Service, sessionService auth.UserTokenService,
	accessControlService accesscontrol.Service,
	apikeyService apikey.Service, userService user.Service,
	jwtService auth.JWTVerifierService, userProtectionService login.UserProtectionService,
	loginAttempts loginattempt.Service, quotaService quota.Service,
	authInfoService login.AuthInfoService, renderService rendering.Service,
	features *featuremgmt.FeatureManager, oauthTokenService oauthtoken.OAuthTokenService,
	socialService social.Service, cache *remotecache.RemoteCache,
	ldapService service.LDAP, settingsProviderService setting.Provider,
	pluginRegistry registry.Service, pluginSettingsService pluginsettings.Service,
) Registration {
	logger := log.New("authn.registration")

	authnSvc.RegisterClient(clients.ProvideRender(renderService))
	authnSvc.RegisterClient(clients.ProvideAPIKey(apikeyService))

	if cfg.LoginCookieName != "" {
		authnSvc.RegisterClient(clients.ProvideSession(cfg, sessionService, authInfoService))
	}

	var proxyClients []authn.ProxyClient
	var passwordClients []authn.PasswordClient
	if cfg.LDAPAuthEnabled {
		ldap := clients.ProvideLDAP(cfg, ldapService, userService, authInfoService)
		proxyClients = append(proxyClients, ldap)
		passwordClients = append(passwordClients, ldap)
	}

	if !cfg.DisableLogin {
		grafana := clients.ProvideGrafana(cfg, userService)
		proxyClients = append(proxyClients, grafana)
		passwordClients = append(passwordClients, grafana)
	}

	// if we have password clients configure check if basic auth or form auth is enabled
	if len(passwordClients) > 0 {
		passwordClient := clients.ProvidePassword(loginAttempts, passwordClients...)
		if cfg.BasicAuthEnabled {
			authnSvc.RegisterClient(clients.ProvideBasic(passwordClient))
		}

		if !cfg.DisableLoginForm {
			authnSvc.RegisterClient(clients.ProvideForm(passwordClient))
		}
	}

	if cfg.AuthProxy.Enabled && len(proxyClients) > 0 {
		proxy, err := clients.ProvideProxy(cfg, cache, proxyClients...)
		if err != nil {
			logger.Error("Failed to configure auth proxy", "err", err)
		} else {
			authnSvc.RegisterClient(proxy)
		}
	}

	if cfg.JWTAuth.Enabled {
		authnSvc.RegisterClient(clients.ProvideJWT(jwtService, cfg))
	}

	if cfg.ExtJWTAuth.Enabled && features.IsEnabledGlobally(featuremgmt.FlagAuthAPIAccessTokenAuth) {
		authnSvc.RegisterClient(clients.ProvideExtendedJWT(cfg))
	}

	for name := range socialService.GetOAuthProviders() {
		clientName := authn.ClientWithPrefix(name)
		authnSvc.RegisterClient(clients.ProvideOAuth(clientName, cfg, oauthTokenService, socialService, settingsProviderService, features))
	}

	// FIXME (jguer): move to User package
	userSync := sync.ProvideUserSync(userService, userProtectionService, authInfoService, quotaService)
	orgSync := sync.ProvideOrgSync(userService, orgService, accessControlService, cfg)
	authnSvc.RegisterPostAuthHook(userSync.SyncUserHook, 10)
	authnSvc.RegisterPostAuthHook(userSync.EnableUserHook, 20)
	authnSvc.RegisterPostAuthHook(orgSync.SyncOrgRolesHook, 30)
	authnSvc.RegisterPostAuthHook(userSync.SyncLastSeenHook, 130)
	authnSvc.RegisterPostAuthHook(sync.ProvideOAuthTokenSync(oauthTokenService, sessionService, socialService).SyncOauthTokenHook, 60)
	authnSvc.RegisterPostAuthHook(userSync.FetchSyncedUserHook, 100)

	rbacSync := sync.ProvideRBACSync(accessControlService)
	if features.IsEnabledGlobally(featuremgmt.FlagCloudRBACRoles) {
		authnSvc.RegisterPostAuthHook(rbacSync.SyncCloudRoles, 110)

		_, exists := pluginRegistry.Plugin(context.Background(), cloudHomePluginId, "")
		if exists {
			authnSvc.RegisterPreLogoutHook(gcomsso.ProvideGComSSOService(pluginSettingsService).LogoutHook, 50)
		}
	}

	authnSvc.RegisterPostAuthHook(rbacSync.SyncPermissionsHook, 120)
	authnSvc.RegisterPostLoginHook(orgSync.SetDefaultOrgHook, 140)

	return Registration{}
}
