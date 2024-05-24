namespace Defyle.WebApi.Policy.Requests
{
  using System.Collections.Generic;
  using System.ComponentModel.DataAnnotations;

  public class DeleteManyPoliciesRequest
  {
    [Required]
    public IEnumerable<string> Ids { get; set; }
  }
}