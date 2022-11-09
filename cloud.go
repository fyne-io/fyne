package fyne

// CloudProvider specifies the identifying information of a cloud provider.
// This information is mostly used by the `fyne.io/cloud ShowSettings' user flow.
//
// Since: 2.3
type CloudProvider interface {
	// ProviderDescription returns a more detailed description of this cloud provider.
	ProviderDescription() string
	// ProviderIcon returns an icon resource that is associated with the given cloud service.
	ProviderIcon() Resource
	// ProviderName returns the name of this cloud provider, usually the name of the service it uses.
	ProviderName() string

	// Cleanup is called when this provider is no longer used and should be disposed.
	// This is guaranteed to execute before a new provider is `Setup`
	Cleanup(App)
	// Setup is called when this provider is being used for the first time.
	// Returning an error will exit the cloud setup process, though it can be retried.
	Setup(App) error
}

// CloudProviderPreferences interface defines the functionality that a cloud provider will include if it is capable
// of synchronizing user preferences.
//
// Since: 2.3
type CloudProviderPreferences interface {
	// CloudPreferences returns a preference provider that will sync values to the cloud this provider uses.
	CloudPreferences(App) Preferences
}

// CloudProviderStorage interface defines the functionality that a cloud provider will include if it is capable
// of synchronizing user documents.
//
// Since: 2.3
type CloudProviderStorage interface {
	// CloudStorage returns a storage provider that will sync documents to the cloud this provider uses.
	CloudStorage(App) Storage
}
