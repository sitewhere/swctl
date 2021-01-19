/**
 * Copyright © 2014-2021 The SiteWhere Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package action

const (
	sitewhereRepoName        = "sitewhere"
	sitewhereRepoURL         = "https://sitewhere.io/helm-charts"
	sitewhereChartName       = "sitewhere-infrastructure"
	sitewhereSystemNamespace = "sitewhere-system"
	sitewhereReleaseName     = "sitewhere"
)

const (
	// ErrIstioNotInstalled is the error when istio is not installed
	ErrIstioNotInstalled = "Istio is not intalled, install istio with `istioctl install` and try again"
)
