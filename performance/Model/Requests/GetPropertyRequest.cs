namespace Defyle.WebApi.Model.Requests
{
  using System.ComponentModel.DataAnnotations;

  public class GetPropertyRequest
  {
    [Required]
    public string Model { get; set; }

    [Required]
    public string Property { get; set; }
  }
}