package gapi

import (
	"context"
	"log"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type Metadata struct {
	UserAgent string
	ClientIP  string
}

const (
	// value yang dihasilkan metadata biasanya adalah map dari dengan key string dan value list of string
	// untuk mengeceknya kita bisa log md dari metadata.FromIncomingcontext
	gprcGatewayUserAgentHeader = "grpcgateway-user-agent" // ini adalah key untuk user agent di metadata pada server gateway
	grpcUserAgentHeader        = "user-agent"             // ini adalah key untuk user agent di metadata pada server grpc
	xForwardedForHeader        = "x-forwarded-for"        // ini adalah key untuk ip address di metadata
)

// kita menggunakan context sebagai input karena semua gRPC metadata akan disimpan di context ini
func (server *Server) extractMetadata(ctx context.Context) *Metadata {
	mtdt := &Metadata{
		UserAgent: "",
		ClientIP:  "",
	}

	// metadata adalah sub package dari grpc yang membantu kita bekerja dengan metadata
	// FromIncomingContext returns the incoming metadata in ctx if it exists.
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		// ok memberitahu kita apakah metadata berhasil didapatkan atau tidak
		// log.Printf("md: %+v\n", md) // print out the content of the md object to see what's inside

		// jika userAgents tidak kosong maka value dari userAgents berada pada item pertama pada list
		// md.get get value of the grpc gateway
		if userAgents := md.Get(gprcGatewayUserAgentHeader); len(userAgents) > 0 {
			mtdt.UserAgent = userAgents[0]
		}

		if userAgents := md.Get(grpcUserAgentHeader); len(userAgents) > 0 {
			mtdt.UserAgent = userAgents[0]
		}

		// ip address juga sama
		if clientIPs := md.Get(xForwardedForHeader); len(clientIPs) > 0 {
			mtdt.ClientIP = clientIPs[0]
		}

	}

	// tidak seperti http, grpc client IP addres tidak di store pada metadata melainkan langsung pada ctx
	// peer juga sub package dari grpc. function ini mengembalikan peer information together with boolean value to tell us if an info exist or not
	if p, ok := peer.FromContext(ctx); ok {
		// jika ok. maka client IPAddress akan tersimpan di field Addr pada p
		mtdt.ClientIP = p.Addr.String()
	}

	log.Print(mtdt)
	return mtdt
}
