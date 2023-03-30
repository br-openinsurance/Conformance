import argparse
from os import listdir
import re


def main(argv = None):
    parser = argparse.ArgumentParser(
        description="Checks if the file names are correct"
    )
    parser.add_argument(
        "apis",
        nargs="+",
        help="Every api to be checked by the script and their respective versions separated by a '_'. Examples: business_1.3, personal_1.0, consents_2.2."
    )
    args = parser.parse_args(argv)
    
    api_list = [api.split('_') for api in args.apis]
    directories = [f"./submissions/functional/{api}/{version}.0" for api, version in api_list]

    wrong_files = []
    for (api, version), directory in zip(api_list, directories):
        regex_pattern = r"^\d{8}_.+_(?P<api>[A-Za-z-]+)_v[12](.[0-9])?(?P<appends>(-[A-Z]{2,4})*)_(0[1-9]|[12]\d|3[01])-(0[1-9]|1[012])-(20\d\d)\."
        if version == '1.0':
            regex_pattern += r"(zip|ZIP)$"
        else:
            regex_pattern += r"json$"

        pattern = re.compile(regex_pattern)

        for file in listdir(directory):
            m = pattern.match(file)
            if (m is None or m.group('api') != api or len(file.split('_')) != 5) and file != ".DS_Store" and file != "readme.md":
                wrong_files.append(file)

    if len(wrong_files):
        print("The following file names are wrong: " + str(wrong_files))
        return 1
        
    return 0


if __name__ == '__main__':
    raise SystemExit(main())
