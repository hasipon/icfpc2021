package ;
import haxe.Json;
import haxe.io.Bytes;
import haxe.macro.Context;
import sys.FileSystem;
import sys.io.File;

class MacroMain
{
	public static function build():Void
	{
		var problems = [];
		for (i in 1...100000)
		{
			if (FileSystem.exists("../problems/" + i))
			{
				problems.push(Json.parse(File.getContent("../problems/" + i)));
			}
			else
			{
				break;
			}
		}
		Context.addResource("problems", Bytes.ofString(Json.stringify(problems)));
	}
}
