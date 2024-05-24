namespace Defyle.WebApi.Auth.Services
{
  using System.Threading.Tasks;
  using AutoMapper;
  using Core.Auth.Models;
  using Core.Auth.Services;
  using Core.Policy.Services;
  using Dtos;

  public class UserDtoService
  {
    private readonly UserService _userService;
    private readonly PolicyService _policyService;
    private readonly IMapper _mapper;

    public UserDtoService(
      UserService userService,
      PolicyService policyService,
      IMapper mapper)
    {
      _userService = userService;
      _policyService = policyService;
      _mapper = mapper;
    }

    public async Task<UserDto> CreateAsync(User user)
    {
      User systemUser = await _userService.FindSystemUserAsync();
      
      UserDto dto = _mapper.Map<UserDto>(user);
      dto.Permissions = await _policyService.GetModelPermissionsForUserAsync(user.Id, "system", systemUser);

      return dto;
    }
  }
}