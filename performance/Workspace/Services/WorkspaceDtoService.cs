namespace Defyle.WebApi.Workspace.Services
{
  using System.Collections.Generic;
  using System.Linq;
  using System.Threading.Tasks;
  using AutoMapper;
  using Core.Auth.Models;
  using Core.Auth.Services;
  using Core.Model.Services;
  using Core.Policy.Services;
  using Core.Workspace.Models;
  using Dtos;

  public class WorkspaceDtoService
	{
    private readonly PolicyService _policyService;
    private readonly UserService _userService;
    private readonly ModelService _modelService;
    private readonly IMapper _mapper;

    public WorkspaceDtoService(
      PolicyService policyService,
      UserService userService,
      ModelService modelService,
      IMapper mapper)
    {
      _policyService = policyService;
      _userService = userService;
      _modelService = modelService;
      _mapper = mapper;
    }

		public async Task<WorkspaceDto> CreateAsync(Workspace workspace, User user)
    {
      WorkspaceDto dto = new WorkspaceDto();

			dto.Id = workspace.Id;
			dto.PartitionId = workspace.PartitionId;
			dto.Name = workspace.Name;
			dto.Image = workspace.Image;
      dto.Encrypted = workspace.Encrypted;
      dto.StorageCapacity = workspace.StorageCapacity;
      dto.CreatedTime = workspace.CreateTime;
      dto.UpdatedTime = workspace.UpdateTime;

      List<WorkspaceTaskDto> tasks = new List<WorkspaceTaskDto>();
			foreach (WorkspaceTask task in workspace.Tasks)
			{
				tasks.Add(_mapper.Map<WorkspaceTaskDto>(task));
			}

			dto.Tasks = tasks;

      User systemUser = await _userService.FindSystemUserAsync();

      if (user.IsSuperuser)
      {
        dto.Permissions = await _modelService.GetPermissionsAsync("workspace");
      }
      else
      {
        var policies = await _policyService.FindAllForUserAsync(user.Id, workspace.Id, systemUser);
        dto.Permissions = policies.Select(e => e.Permission);
      }

      return dto;
		}
	}
}