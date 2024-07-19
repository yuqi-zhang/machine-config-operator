package build

// Label that associates any objects with on-cluster layering. Should be added
// to every object that BuildController creates or manages, ephemeral or not.
const (
	OnClusterLayeringLabelKey = "machineconfiguration.openshift.io/on-cluster-layering"
)

// Labels to add to all ephemeral build objects the BuildController creates.
const (
	EphemeralBuildObjectLabelKey    = "machineconfiguration.openshift.io/ephemeral-build-object"
	RenderedMachineConfigLabelKey   = "machineconfiguration.openshift.io/rendered-machine-config"
	TargetMachineConfigPoolLabelKey = "machineconfiguration.openshift.io/target-machine-config-pool"
)

// Annotations to add to all ephemeral build objects BuildController creates.
const (
	machineOSBuildNameAnnotationKey  = "machineconfiguration.openshift.io/machine-os-build"
	machineOSConfigNameAnnotationKey = "machineconfiguration.openshift.io/machine-os-config"
)

// Entitled build secret names
const (
	// Name of the etc-pki-entitlement secret from the openshift-config-managed namespace.
	EtcPkiEntitlementSecretName = "etc-pki-entitlement"

	// Name of the etc-pki-rpm-gpg secret.
	EtcPkiRpmGpgSecretName = "etc-pki-rpm-gpg"

	// Name of the etc-yum-repos-d ConfigMap.
	EtcYumReposDConfigMapName = "etc-yum-repos-d"
)

// Canonical secrets
const (
	canonicalSecretSuffix string = "-canonical"
	// This label gets applied to all secrets that we've canonicalized as a way
	// to indicate that we created and own them.
	CanonicalSecretLabelKey string = "machineconfiguration.openshift.io/canonicalizedSecret"
	// This label is applied to all canonicalized secrets. Its value should
	// contain the original name of the secret that has been canonicalized.
	OriginalSecretNameLabelKey string = "machineconfiguration.openshift.io/originalSecretName"
)

const (
	// Filename for the machineconfig JSON tarball expected by the build pod
	machineConfigJSONFilename string = "machineconfig.json.gz"
)

// Entitled build annotation keys
const (
	entitlementsAnnotationKeyBase  = "machineconfiguration.openshift.io/has-"
	EtcPkiEntitlementAnnotationKey = entitlementsAnnotationKeyBase + EtcPkiEntitlementSecretName
	EtcYumReposDAnnotationKey      = entitlementsAnnotationKeyBase + EtcYumReposDConfigMapName
	EtcPkiRpmGpgAnnotationKey      = entitlementsAnnotationKeyBase + EtcPkiRpmGpgSecretName
)
