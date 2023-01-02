package java

import (
	"dagger.io/dagger"
)

// WithMaven returns a dagger container with sane defaults
// for building java maven apps using the Eclipse temurin JDK
func WithMaven(c *dagger.Client) *dagger.Container {
	hostDir := c.Host().Directory(".", dagger.HostDirectoryOpts{
		// exclude target to avoid cache invalidation
		Exclude: []string{"target"},
	})

	// TODO this cache won't work in CI's were the docker
	// engine is ephemeral (i.e GHA). We should do optimizations
	// for those cases
	mvnCache := c.CacheVolume("mvn-cache")

	return c.Container().From("maven:3.8.6-eclipse-temurin-17").
		// TODO what can we do to speed this up in a CI so we don't
		// have to pull the dependencies each time?
		WithMountedCache("/root/.m2", mvnCache).
		WithMountedDirectory("/app", hostDir).
		WithWorkdir("/app")
}
