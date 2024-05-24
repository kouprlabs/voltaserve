# appsettings.json
To enable LDAP authentication, `AuthenticationType` must be set to `ldap`, for local authentication, the value `local` must be set instead.

`Ldap.BindDn` can contain values like `DEFYLE\\WebService` or a full Distinguished Name like `CN=WebService,CN=Users,DC=Defyle,DC=local`.

To get the Distinguished Names of users, run `dsquery user` on a Windows Server Command Prompt.

`SecurityKey` must be larger than or equal to 256 bytes.
