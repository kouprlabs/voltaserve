namespace Defyle.WebApi.Inode.Controllers
{
  using System.Threading.Tasks;
  using Core.Auth.Services;
  using Core.Inode.Services;
  using Core.Workspace.Services;
  using Infrastructure.Controllers;

  public class BaseInodeController : BaseController
  {
    private readonly InodeService _service;
    private readonly WorkspaceService _workspaceService;
    private readonly UserService _userService;

    public BaseInodeController(
      InodeService service,
      WorkspaceService workspaceService,
      UserService userService)
    {
      _service = service;
      _workspaceService = workspaceService;
      _userService = userService;
    }

    protected async Task<string> GetEffectiveNodeIdAsync(string workspaceId, string id)
    {
      if (id == "0" || string.IsNullOrWhiteSpace(id))
      {
        var user = await _userService.FindAsync(UserId);
        var workspace = await _workspaceService.FindAsync(workspaceId, user);
        var inode = await _service.FindRootAsync(workspace, user);
        return inode.Id;
      }
      else
      {
        return id;
      }
    }
  }
}