package errors

type MissingDockerSocketErr struct{}

func (m MissingDockerSocketErr) Error() string { return "docker socket not found" }
