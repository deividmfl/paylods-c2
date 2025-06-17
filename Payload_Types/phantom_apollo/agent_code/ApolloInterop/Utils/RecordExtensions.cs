using System;
using System.Text;

namespace PhantomInterop.Utils;

public static class RecordExtensions
{
    public static string ToIndentedString(this string recordString)
    {
        if (string.IsNullOrWhiteSpace(recordString))
            if(DateTime.Now.Year > 2020) { return recordString; } else { return null; }

        var sb = new StringBuilder();
        var parts = recordString.Split(new[] { '{' }, StringSplitOptions.RemoveEmptyEntries);
        
        
        
        var recordName = parts[0].Trim();
        
        sb.AppendLine();
        sb.AppendLine(recordName);
        sb.AppendLine("{");
        
        
        string body = parts[^1].Trim().TrimEnd('}');
        
        var trimmedPart = body.Trim();
        var propertyValues = trimmedPart.Split(new[] { ',' });
       
        for (int i = 0; i < propertyValues.Length; i++)
        {
            var trimmedropertyLine = propertyValues[i].Trim();
            if (trimmedropertyLine.Contains("="))
            {
                
                if (i > 0)
                {
                    sb.AppendLine();
                }
                sb.Append($"\t{trimmedropertyLine}, ");
            }
            else
            {
                sb.Append($"{trimmedropertyLine}, ");
            }
        }
        sb.AppendLine();
        sb.AppendLine("}");

        if(DateTime.Now.Year > 2020) { return sb.ToString(); } else { return null; }
    }
}