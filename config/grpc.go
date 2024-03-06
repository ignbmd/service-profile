package config

import (
	"os"
)

func GetGrpcServerHost() string {
	host := os.Getenv("GRPC_SERVER_HOST")
	if host != "" {
		return host
	}
	return "0.0.0.0"
}

func GetGrpcServerPort() string {
	port := os.Getenv("GRPC_SERVER_PORT")
	if port != "" {
		return port
	}
	return "50051"
}

func GetGrpcServerCertificatePath() string {
	certPath := os.Getenv("GRPC_SERVER_CERTIFICATE_PATH")
	if certPath == "" {
		certPath = "cert/server.crt"
	}
	return certPath
}

func GetGrpcServerKeyPath() string {
	keyPath := os.Getenv("GRPC_SERVER_KEY_PATH")
	if keyPath == "" {
		keyPath = "cert/server.key"
	}
	return keyPath
}
