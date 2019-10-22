/*
 * Copyright (c) 2016-2019, Fraunhofer AISEC. All rights reserved.
 *
 *
 *            $$\                           $$\ $$\   $$\
 *            $$ |                          $$ |\__|  $$ |
 *   $$$$$$$\ $$ | $$$$$$\  $$\   $$\  $$$$$$$ |$$\ $$$$$$\    $$$$$$\   $$$$$$\
 *  $$  _____|$$ |$$  __$$\ $$ |  $$ |$$  __$$ |$$ |\_$$  _|  $$  __$$\ $$  __$$\
 *  $$ /      $$ |$$ /  $$ |$$ |  $$ |$$ /  $$ |$$ |  $$ |    $$ /  $$ |$$ |  \__|
 *  $$ |      $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$ |  $$ |$$\ $$ |  $$ |$$ |
 *  \$$$$$$\  $$ |\$$$$$   |\$$$$$   |\$$$$$$  |$$ |  \$$$   |\$$$$$   |$$ |
 *   \_______|\__| \______/  \______/  \_______|\__|   \____/  \______/ \__|
 *
 * This file is part of Clouditor Community Edition.
 *
 * Clouditor Community Edition is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Clouditor Community Edition is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * long with Clouditor Community Edition.  If not, see <https://www.gnu.org/licenses/>
 */

package io.clouditor.discovery.aws;

import io.clouditor.discovery.ScannerInfo;
import java.util.Collections;
import java.util.List;
import software.amazon.awssdk.services.iam.model.PasswordPolicy;

@ScannerInfo(assetType = "PasswordPolicy", group = "AWS", service = "IAM")
public class AwsIamPasswordPolicyScanner extends AwsIamScanner<PasswordPolicy> {

  /**
   * This is not a really ARN but rather a simple hack to give the password policy an asset
   * identifier.
   */
  static final String ARN_AWS_IAM_PW_POLICY = "arn:aws:iam:::password-policy";

  public AwsIamPasswordPolicyScanner() {
    super(policy -> ARN_AWS_IAM_PW_POLICY, policy -> "Password Policy");
  }

  @Override
  protected List<PasswordPolicy> list() {
    return Collections.singletonList(this.api.getAccountPasswordPolicy().passwordPolicy());
  }
}
