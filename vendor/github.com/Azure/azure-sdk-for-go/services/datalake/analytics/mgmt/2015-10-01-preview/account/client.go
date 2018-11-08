// Package account implements the Azure ARM Account service API version 2015-10-01-preview.
//
// Creates an Azure Data Lake Analytics account management client.
package account

// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.
//
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

import (
	"github.com/Azure/go-autorest/autorest"
)

const (
	// DefaultBaseURI is the default URI used for the service Account
	DefaultBaseURI = "https://management.azure.com"
)

// Deprecated: Please use package github.com/Azure/azure-sdk-for-go/services/preview/datalake/analytics/mgmt/2015-10-01-preview/account instead.
// BaseClient is the base client for Account.
type BaseClient struct {
	autorest.Client
	BaseURI        string
	SubscriptionID string
}

// Deprecated: Please use package github.com/Azure/azure-sdk-for-go/services/preview/datalake/analytics/mgmt/2015-10-01-preview/account instead.
// New creates an instance of the BaseClient client.
func New(subscriptionID string) BaseClient {
	return NewWithBaseURI(DefaultBaseURI, subscriptionID)
}

// Deprecated: Please use package github.com/Azure/azure-sdk-for-go/services/preview/datalake/analytics/mgmt/2015-10-01-preview/account instead.
// NewWithBaseURI creates an instance of the BaseClient client.
func NewWithBaseURI(baseURI string, subscriptionID string) BaseClient {
	return BaseClient{
		Client:         autorest.NewClientWithUserAgent(UserAgent()),
		BaseURI:        baseURI,
		SubscriptionID: subscriptionID,
	}
}
