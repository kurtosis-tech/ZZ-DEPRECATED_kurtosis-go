/*
 * Copyright (c) 2020 - present Kurtosis Technologies LLC.
 * All Rights Reserved.
 */

package services

import "os"

// GENERIC TOOD: If Go had generics, this would be parameterized with the subtype of Service that this returns
type DockerContainerInitializer interface {
	// Gets the Docker image that will be used for instantiating the Docker container
	GetDockerImage() string

	// Gets the "set" of ports that the Docker container running the service will listen on
	GetUsedPorts() map[int]bool

	// GENERICS TOOD: When Go has generics, make this return type be parameterized
	/*
		Uses the IP address of the Docker container running the service to create an implementation of the interface the developer
		has created to represent their service.

		NOTE: Because Go doesn't have generics, we can't properly parameterize the return type to be the actual service interface
		that the developer has created; nonetheless, the developer should return an implementation of their interface (which itself
		should extend Service).

		Args:
			ipAddr: The IP address of the Docker container running the service
	*/
	GetServiceFromIp(ipAddr string) Service

	// GENERICS TOOD: If Go had generics, we could parameterize this entire class with an enum of the types of files this service consumes
	/*
		This method is used to declare that the service will need a set of files in order to run. To do this, the developer
		declares a set of string keys that are meaningful to the developer, and Kurtosis will create one file per key. These newly-createed
		file objects will then be passed in to the `InitializeMountedFiles` and `GetStartCommand` functions below keyed on the
		strings that the developer passed in, so that the developer can initialize the contents of the files as they please.
		Kurtosis then guarantees that these files will be made available to the service at startup time.

		NOTE: The keys that the developer returns here are ONLY used for developer identification purposes; the actual
		filenames and filepaths of the file are implementation details handled by Kurtosis!

		Returns:
			A "set" of user-defined key strings identifying the files that the service will need, which is how files will be
				identified in `InitializeMountedFiles` and `GetStartCommand`
	*/
	GetFilesToMount() map[string]bool

	/*
		Initializes the contents of the files that the developer requested in `GetFilesToMount` with whatever
			contents the developer desires. This will be called before service startup.

		Args:
			mountedFiles: A mapping of developer_key -> file_pointer, with developer_key corresponding to the keys declares in
				`GetFilesToMount`
	*/
	InitializeMountedFiles(mountedFiles map[string]*os.File) error

	/*
		Kurtosis mounts the files that the developer requested in `GetFilesToMount` via a Docker volume, but Kurtosis doesn't
		know anything about the Docker image backing the service so therefore doesn't know what filepath it can safely mount
		the volume on. This function uses the developer's knowledge of the Docker image running the service to inform
		Kurtosis of a filepath where the Docker volume can be safely mounted.

		Returns:
			A filepath on the Docker image backing this service that's safe to mount the test volume on
	*/
	GetTestVolumeMountpoint() string

	/*
		Uses the given arguments to build the command that the Docker container running this service will be launched with.

		NOTE: Because the IP address of the container is an implementation detail, any references to the IP address of the
			container should use the placeholder "SERVICEIP" instead. This will get replaced at launch time with the service's
			actual IP.

		Args:
			mountedFileFilepaths: Mapping of developer_key -> initialized_file_filepath where developer_key corresponds to the keys returned
				in the `GetFilesToMount` function, and initialized_file_filepath is the path *on the Docker container* of where the
				file has been mounted. The files will have already been initialized via the `InitializeMountedFiles` function.
			ipPlaceholder: Because the IP address of the container is an implementation detail, any references to the IP
				address of the container should use this placeholder string when they want the IP. This will get replaced at
				service launch time with the actual IP.

		Returns:
			The command fragments which will be used to construct the run command which will be used to launch the Docker container
				running the service. If this is nil, then no explicit command will be specified and whatever command the Dockerfile
				specifies will be run instead.
	*/
	GetStartCommand(mountedFileFilepaths map[string]string, ipPlaceholder string) ([]string, error)
}
