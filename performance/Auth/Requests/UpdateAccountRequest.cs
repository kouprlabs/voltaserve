namespace Defyle.WebApi.Auth.Requests
{
  using System.ComponentModel.DataAnnotations;

  public class UpdateAccountRequest
  {
    [EmailAddress]
    public string NewEmail { get; set; }

    public string CurrentPassword { get; set; }

    public string NewPassword { get; set; }

    public string NewFullName { get; set; }
  }
}