namespace Defyle.Core.Auth.Services
{
  using System.Drawing;
  using System.Linq;
  using System.Threading.Tasks;
  using AutoMapper;
  using Infrastructure.Poco;
  using Models;
  using Novell.Directory.Ldap;
  using Preview.Services;
  using Proto;

  public class LdapService
  {
    private readonly Size _profileImageSize = new Size(230, 230);
    private readonly UserService _userService;
    private readonly CoreSettings _coreSettings;
    private readonly UserServiceProto.UserServiceProtoClient _protoClient;
    private readonly IMapper _mapper;

    public LdapService(
      UserService userService,
      CoreSettings coreSettings,
      UserServiceProto.UserServiceProtoClient protoClient,
      IMapper mapper)
    {
      _userService = userService;
      _coreSettings = coreSettings;
      _protoClient = protoClient;
      _mapper = mapper;
    }
    
    public async Task<User> CreateFromLdapEntryAsync(LdapEntry ldapEntry)
    {
      string username = ldapEntry.GetAttribute(_coreSettings.Ldap.UsernameAttribute).StringValue;

      User systemUser = await _userService.FindSystemUserAsync();
      User user = await _userService.FindByUsernameAsync(username, systemUser);

      if (user == null)
      {
        user = new User();
        user.IsLdap = true;
        user.Username = username.Trim().ToLowerInvariant();

        UpdateLdapUser(user, ldapEntry, _coreSettings.Ldap.AdminCn);
        
        await _protoClient.InsertAsync(new UserInsertRequestProto
        {
          UserId = systemUser.Id,
          User = _mapper.Map<UserProto>(user)
        });
      }
      else
      {
        UpdateLdapUser(user, ldapEntry, _coreSettings.Ldap.AdminCn);
        
        await _protoClient.UpdateAsync(new UserUpdateRequestProto
        {
          UserId = systemUser.Id,
          User = _mapper.Map<UserProto>(user)
        });
      }

      return user;
    }
    
    private void UpdateLdapUser(User user, LdapEntry ldapEntry, string adminCn)
    {
      LdapAttribute memberOfAttribute = ldapEntry.GetAttribute("memberOf");

      if (memberOfAttribute != null)
      {
        user.IsSuperuser = memberOfAttribute.StringValueArray.Contains(adminCn);
      }

      LdapAttribute displayNameAttribute = ldapEntry.GetAttribute("displayName");

      if (displayNameAttribute != null)
      {
        user.FullName = displayNameAttribute.StringValue;
      }

      LdapAttribute mailAttribute = ldapEntry.GetAttribute("mail");

      if (mailAttribute != null)
      {
        string mail = mailAttribute.StringValue;

        if (!Base64ImageService.IsValidBase64Image(user.Image) &&
            mail != user.Email &&
            _coreSettings.GravatarIntegration)
        {
          user.Image = Base64ImageService.GenerateGravatar(mail, _profileImageSize.Width);
        }

        user.Email = mail.Trim().ToLowerInvariant();
      }
    }
  }
}