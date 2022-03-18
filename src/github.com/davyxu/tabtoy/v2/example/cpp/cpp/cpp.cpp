// cpp.cpp : �������̨Ӧ�ó������ڵ㡣
//

#include "stdafx.h"
#include "Logger.h"
#include "DataReader.h"
#include "Config.h"

int main()
{
	std::ifstream stream("../../csharp/Example/Config.bin", std::ios::binary);
	tabtoy::DataReader reader(stream);
	if (!reader.ReadHeader())
	{
		//Console.WriteLine("combine file crack!");
		return 1;
	}

	table::Config config;
	table::Config::Deserialize( config, reader);
		

    return 0;
}

