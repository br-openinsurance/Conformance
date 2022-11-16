import argparse
from os import listdir
import re

def main(argv = None):
    parser = argparse.ArgumentParser(
        description='Checks the zip files added are correct'
    )
    parser.add_argument(
        "version",
        help="Api version."
    )
    parser.add_argument(
        "apis",
        nargs="+",
        help="Every api to be checked by the script."
    )
    args = parser.parse_args(argv)
    
    api_list = args.apis
    version  = args.version
    directories = [f"./submissions/functional/{api}/{version}" for api in api_list]

    wrong_zips = []
    for api, directory in zip(api_list, directories):
        pattern = re.compile(r"^\d{8}_.+_(?P<api>[A-Za-z-]+)_v[12](-OL)?_(0[1-9]|[12]\d|3[01])-(0[1-9]|1[012])-(20\d\d).(zip|ZIP)$")
        for file in listdir(directory):
            m = pattern.match(file)
            if (m is None or (m.group('api') != api)) and file != ".DS_Store":
                wrong_zips.append(file)

    if len(wrong_zips):
        print("The following zip names are wrong: " + str(wrong_zips))
        return 1
        
    return 0


if __name__ == '__main__':
    raise SystemExit(main())
