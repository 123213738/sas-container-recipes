## Contents

- [Configure Your Environment with SAS Environment Manager](#configure-your-environment-with-sas-environment-manager)
- [Verify That Licenses Are Applied](#verify-that-licenses-are-applied)
- [(Optional) Create a Local Copy of Documentation](#optional-create-a-local-copy-of-documentation)

## Configure Your Environment with SAS Environment Manager

**Note:** The tasks in this section are applicable for a full deployment. Skip this section if you performed a programming-only deployment.

### Sign In as the sasboot User

Your SAS environment is deployed with an initial administrator account that is named sasboot. The password for this account has expired by default, so you must reset the password before you can sign in.

To reset the password:

1. Locate the most recent log for the SAS Logon service in /var/log/sas/viya/saslogon/default.

1. Search the log for the characters sasboot:

    ```
    kubectl -n <k8s namespace> get pods | grep coreservices
    kubectl -n <k8s namespace> exec -it sas-viya-coreservices-<uuid> -- ls -l /var/log/sas/viya/saslogon/default
    kubectl -n <k8s namespace> exec -it sas-viya-coreservices-<uuid> -- grep 'sasboot' /var/log/sas/viya/saslogon/default/sas- 
    saslogon_date-and-time-stamp.log 2> /dev/null
    ```
    Here is a typical message:

    ```
    Reset password for initial user sasboot using link: /SASLogon/reset_password?code=xxxxxx
    ```

3. Sign in from a URL with this format:

    `https://igress-path/SASLogon/reset_password?code=password`

    Make a note of this URL to share with any other users of your SAS Viya software.

4. Follow the instructions on the displayed web page to reset the password.

    **Note:** If the URL has expired, you will need to restart the coreservices pod:

    `kubectl -n <k8s namespace> delete pod sas-viya-coreservices-<uuid>`

    Then go to the log and obtain the new URL. The URL expires 24 hours after the SAS Logon service restarts. For security purposes, the URL that is specified in a browser or in a text editor also expires, even if the password is not reset.

    After you reset the password, SAS Environment Manager automatically opens in your browser.

5. Click Yes for all the assumable groups so that you have the permissions to perform subsequent tasks.

### Configure the Connection to Your Identity Provider

Configure the connection to your identity provider before your users can access SAS Environment Manager and SAS Visual Analytics.

While signed in as sasboot, configure the connection to your identity provider:

**Note:** Only LDAP-based identity providers are supported. You need to have basic familiarity with LDAP administration. For more information about the properties that are relevant for this procedure, see [sas.identities.providers.ldap](http://documentation.sas.com/?cdcId=calcdc&amp;cdcVersion=3.4&amp;docsetId=calconfig&amp;docsetTarget=n08000sasconfiguration0admin.htm#n08044sasconfiguration0admin) in  _SAS Viya Administration: Configuration Properties_.

1. Select the wrench from the side menu to open the Configuration page.

1. On the Configuration page, select **Basic Services** from the list, and then select the **Identities** service from the list of services.

3. To configure user properties, in the **sas.identities.providers.ldap.user** section, click **New Configuration**:

    1. Specify a value for the **baseDN** required field. For the remaining fields, review the default values and make changes, as necessary. The default values are appropriate for most sites.

        **Note:** When using the LDAP protocol, passwords are transmitted over the network as clear-text. To secure the deployment, SAS recommends that you configure encrypted LDAP connections. For more information, see [Encrypt LDAP Connections](http://documentation.sas.com/?docsetId=calencryptmotion&amp;docsetTarget=n1xdqv1sezyrahn17erzcunxwix9.htm&amp;docsetVersion=3.4#p1bai319815977n1bzdyyxr3d5he) in _Encryption in SAS Viya: Data in Motion_.

        For each property that represents a user-level field in SAS, specify a corresponding property in the LDAP server software.

        **Tip:** Consider specifying a custom filter to limit the user accounts that SAS Viya returns from your LDAP server.
 
    2. Click **Save**.

4. To configure group properties, in the **sas.identities.providers.ldap.group** section, click **New Configuration**. In the **New Configuration** window:

    1. Specify a value for the **baseDN** required field. For the remaining fields, review the default values and make changes, as necessary. The default values are appropriate for most sites.

        For each property that represents a group-level field in SAS, specify a corresponding property in the LDAP server software.

        **Tip:** Consider specifying a custom filter to limit the accounts that SAS Viya returns from your LDAP server.

    2. Click **Save**.

5. To configure connection properties, in the **sas.identities.providers.ldap.connection** section, click **New Configuration**. In the **New Configuration** window:

    1. Specify values for the following required fields: **host** , **password** , **port** , **url** , and **userDN**. For the remaining fields, review the default values and make changes, as necessary. The default values are appropriate for most sites.

    2. Click **Save**.

6. To verify user and group information, from the SAS Environment Manager side menu, select **New** to open the Users page.

    1. On the Users page, select **Users** from the list in the toolbar. Your users should appear after a few minutes. It is not necessary to restart any servers or services. Then select **Groups** from the list to display your groups.

    2. Verify that user and group information is displayed correctly. If not, make any necessary changes to the identities service properties.

7. Restart the Identities and SAS Logon Manager services by starting a new coreservices pod:

    `kubectl -n <k8s namespace> delete pod sas-viya-coreservices-<uuid>`

### Set Up Administrative Users

While you are signed on to SAS Environment Manager as the sasboot user, set up at least one SAS Administrator user, as follows:

1. On the Users page in SAS Environment Manager, select **Custom Groups** from the list in the toolbar.

1. In the left pane, click **SAS Administrators**.

1. In the **Members** section of the right pane, click the **Edit** icon, and add one or more members to the group (including your own account, if applicable).

1. Sign out from SAS Environment Manager so that you are no longer signed in as the sasboot user.

1. If you added your own account to the SAS Administrators group, you can sign on again to SAS Environment Manager using that account.

    Open SAS EnvironmentManager from a URL with the following format:

    `https://ingress-path/SASEnvironmentManager`

    **Tip:** Since SAS Administrators is an assumable group, the following prompt is displayed: **Do you want to opt in to all of your assumable groups?**. Select **Yes** if you want the extra permissions that are associated with the SAS Administrators group. The selection remains in effect until you sign out.

### Sign In Using LDAP Credentials

Open SAS Environment Manager from a URL with the following format:

`https://ingress-path/SASEnvironmentManager`

Sign in as one of the SAS Administrators that you set up in the Set Up Administrative Users section above.

### Configure the Connection to the Mail Service

Configure the connection to your mailservice. Complete these steps while you are signed in as one of theSAS Administrators.

1. Select the wrench icon from the side menu to open the Configuration page.

1. On the Configuration page, select **Basic Services** from the list, and then select **Mail service** from the list of services.

3. In the **sas.mail** section, click the Edit icon. In the Edit Configuration window,follow these steps:

    1. Specify a value for the following required fields: **host** and **port**. For the remaining fields, review the default values and make changes, as necessary. The default values are appropriate for most sites.

    1. Click **Save**.

4. (Optional) To enable the health check for the mail service, perform the following steps.

    1. Select the wrench icon from the side menu to open the Configuration page.

    1. On the Configuration page, select **Basic Services** from the list, and then select **Mail service** from the list of services.

    1. In the **management.health.mail** section, click the Edit icon.

    1. Turn the **enabled** toggle to **on**.

    1. Click **Save**.

    When this toggle is set, health checks will be enabled after the mail service is restarted. If the mail host is not configured or is configured incorrectly, or if it cannot connect to the SMTP mail server, the mail service will indicate it is in a failed state.

5. Restart the mail service.

    `kubectl -n <k8s namespace> delete pod sas-viya-coreservices-<uuid>`

### Disable the Password Reset Feature and Reset the sasboot Password

When you are finished setting up LDAP and the initial administrative users, you should reset the password for the sasboot user. For additional security, you can then disable the password reset feature. This action prevents password reset links from being written to the log each time the SASLogon service is restarted.

1. Sign in to SAS Environment Manager as an administrative user and select the wrench icon from the side menu to open the Configuration page.

1. On the Configuration page, select **Definitions** from the drop-down list.

1. In the left pane, select **sas.logon.initial**. Then click **New Configuration** at the top of the right pane. If a definition already exists, you can select the Edit icon to edit the existing definition.

1. In the New sas.logon.initial Configuration window or the Edit sas.logon.initial Configuration window, set **reset.enabled** to **off**.

1. Click **Save**.

1. Restart the SAS LogonManager services:

    `kubectl -n <k8s namespace> delete pod sas-viya-coreservices-<uuid>`

    For more information, see [General Servers and Services: Operate](http://documentation.sas.com/?cdcId=calcdc&amp;cdcVersion=3.4&amp;docsetId=calchkadm&amp;docsetTarget=n00003ongoingtasks00000admin.htm&amp;locale=en) in  _SAS Viya Administration: General Servers and Services_.

    **Note:** After you disable this feature, you can still change the sasboot password if the existing password is known. Enter the URL for SAS Viya with the path /SASLogon/change_password. If you are already signed in as another user, first sign out and then sign back in as sasboot using the current password. You can then complete the steps to change the password.

### Configure SAS Viya to Encrypt the LDAP Connection

SAS Viya supports encrypted connections between the LDAP client and server. To configure a secure LDAP connection, see [Encrypt  LDAP Connections](http://documentation.sas.com/?docsetId=calencryptmotion&amp;docsetTarget=n1xdqv1sezyrahn17erzcunxwix9.htm&amp;docsetVersion=3.4#p1bai319815977n1bzdyyxr3d5he) in _Encryption in SAS Viya: Data in Motion._

## Verify That Licenses Are Applied

During deployment, a license is applied to both the CAS in-memory compute engine and the SAS Foundation compute engine. To ensure proper operation of the engines, you should verify that the licenses were applied properly.

For more information, see [Licensing: How To (SAS Studio)](https://go.documentation.sas.com/?cdcId=calcdc&cdcVersion=3.4&docsetId=callicense&docsetTarget=n03028saslicensing00000admin.htm) in _SAS Viya Administration_.

If the licenses were not applied, use the instructions to apply the licenses.

## (Optional) Create a Local Copy of Documentation

You can configure your software to give your users access to local documentation. Here are two instances where access to local documentation would be useful:

- You have customized your documentation.
- Your SAS system is highly secure, and it does not have access to the internet. Because the SAS documentation is cloud-hosted, it cannot be reached without internet access.

**Note:** The cloud-hosted SAS documentation is frequently updated. The SAS administrator should refresh the local copy on a regular basis to ensure that your users have up-to-date information.

You can download PDF versions of the documentation, or you can create customized versions of the documentation. Create an HTML page with links to all documents that make up your local documentation collection, and create a link to this page.

To configure local documentation:

1. Access SAS Environment Manager.

1. Select **Configuration** from the left navigation bar.

1. Under the **View** menu, select **Definitions**.

1. Select the sas.htmlcommons definition.

1. Click **New Configuration**.

1. On the New sas.htmlcommons Configuration pane, click **Add Property** to add the following two properties:
    - **additionalHelpMenuUrl** — Specify the path to the HTML page that contains links to your local documentation.
    - **additionalHelpMenuLabel** — Provide a meaningful label for the link that your users can access.
      If you do not provide this parameter, a default label of **Additional Help** is used.

1. Click **Save** on the New htmlcommons Configuration pane.

Users see a new item in the **Help Menu** list, between the **Help Center** and **About** entries. Clicking this link opens the specified HTML page.