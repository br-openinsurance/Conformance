import argparse
import os
import re
import json
from .load_allow_list import load_allow_list

ALLOWLIST = load_allow_list()


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


def is_invalid_filename(filename, api, version):
    if filename in [".DS_Store", "readme.md"]:
        return True

    regex_pattern = r"^\d{8}_.+_(?P<api>[A-Za-z-]+)_v[12](\.[0-9][0-9]?)?(?P<appends>(-[A-Z]{2,4})*)_(0[1-9]|[12]\d|3[01])-(0[1-9]|1[012])-(20\d\d)\."

    if version == '1.0':
        regex_pattern += r"(zip|ZIP)$"
    else:
        regex_pattern += r"json$"

    pattern = re.compile(regex_pattern)
    m = pattern.match(filename)

    if m is None:
        print(f"Regex match failed: {filename}")
        return True

    if m.group('api') != api:
        print(f"API mismatch in filename '{filename}': expected '{api}', got '{m.group('api')}'")
        return True

    if len(filename.split('_')) != 5:
        print(f"Filename split issue: expected 5 parts separated by '_', got {len(filename.split('_'))} — {filename}")
        return True

    return False


def check_filenames(apis):
    wrong_files = {}

    for api in apis:
        api_name, version = api.split('_')
        directory = f"./submissions/functional/{api_name}/{version}.0"
        wrong_files[directory] = []
        
        for file in os.listdir(directory):
            if os.path.basename(file) in ALLOWLIST:
                continue
            if is_invalid_filename(file, api_name, version):
                wrong_files[directory].append(file)
        # sort the array of wrong files for better readability
        wrong_files[directory].sort()

    return {dir:files for dir, files in wrong_files.items() if files}


def main():
    args = parse_args()
    wrong_files = check_filenames(args.apis)

    if wrong_files:
        print("The following file names are wrong:\n" + json.dumps(wrong_files, indent=4, ensure_ascii=False))
        return 1

    return 0


if __name__ == '__main__':
    raise SystemExit(main())