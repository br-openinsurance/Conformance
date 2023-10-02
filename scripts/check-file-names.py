import argparse
import os
import re


def parse_args():
    parser = argparse.ArgumentParser(
        description="Checks if the file names are correct"
    )
    parser.add_argument(
        "apis",
        nargs="+",
        help="Every api to be checked by the script and their respective versions separated by a '_'. Examples: business_1.3, personal_1.0, consents_2.2."
    )
    return parser.parse_args()


def is_valid_filename(filename, api, version):
    regex_pattern = r"^\d{8}_.+_(?P<api>[A-Za-z-]+)_v[12](.[0-9])?(?P<appends>(-[A-Z]{2,4})*)_(0[1-9]|[12]\d|3[01])-(0[1-9]|1[012])-(20\d\d)\."

    if version == '1.0':
        regex_pattern += r"(zip|ZIP)$"
    else:
        regex_pattern += r"json$"

    pattern = re.compile(regex_pattern)
    m = pattern.match(filename)
    
    return (m is None or m.group('api') != api or len(filename.split('_')) != 5) and filename != ".DS_Store" and filename != "readme.md"


def check_filenames(apis):
    wrong_files = []

    for api in apis:
        api_name, version = api.split('_')
        directory = f"./submissions/functional/{api_name}/{version}.0"

        for file in os.listdir(directory):
            if is_valid_filename(file, api_name, version):
                wrong_files.append(file)

    return wrong_files


def main():
    args = parse_args()
    wrong_files = check_filenames(args.apis)

    if wrong_files:
        print("The following file names are wrong: " + str(wrong_files))
        return 1

    return 0


if __name__ == '__main__':
    raise SystemExit(main())
