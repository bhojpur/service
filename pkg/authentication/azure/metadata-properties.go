package azure

// Copyright (c) 2018 Bhojpur Consulting Private Limited, India. All rights reserved.

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

// MetadataKeys : Keys for all metadata properties.
var MetadataKeys = map[string][]string{ // nolint: gochecknoglobals
	// clientId, clientSecret, tenantId are supported for backwards-compatibility as they're used by some components, but should be considered deprecated

	// Certificate contains the raw certificate data
	"Certificate": {"azureCertificate", "spnCertificate"},
	// Path to a certificate
	"CertificateFile": {"azureCertificateFile", "spnCertificateFile"},
	// Password for the certificate
	"CertificatePassword": {"azureCertificatePassword", "spnCertificatePassword"},
	// Client ID for the Service Principal
	// The "clientId" alias is supported for backwards-compatibility as it's used by some components, but should be considered deprecated
	"ClientID": {"azureClientId", "spnClientId", "clientId"},
	// Client secret for the Service Principal
	// The "clientSecret" alias is supported for backwards-compatibility as it's used by some components, but should be considered deprecated
	"ClientSecret": {"azureClientSecret", "spnClientSecret", "clientSecret"},
	// Tenant ID for the Service Principal
	// The "tenantId" alias is supported for backwards-compatibility as it's used by some components, but should be considered deprecated
	"TenantID": {"azureTenantId", "spnTenantId", "tenantId"},
	// Identifier for the Azure environment
	// Allowed values (case-insensitive): AZUREPUBLICCLOUD, AZURECHINACLOUD, AZUREGERMANCLOUD, AZUREUSGOVERNMENTCLOUD
	"AzureEnvironment": {"azureEnvironment"},
}

// Default Azure environment.
const DefaultAzureEnvironment = "AZUREPUBLICCLOUD"
