from io import StringIO, TextIOBase
import logging
from sys import argv
from typing import override


logger = logging.getLogger(__name__)


class database:
    def __init__(self, valid_ranges: list[tuple[int, int]] | None = None) -> None:
        self.valid_ranges: list[tuple[int, int]] = valid_ranges or []

    @override
    def __repr__(self) -> str:
        output = "database:\n"

        index = 0
        for r in self.valid_ranges:
            output += f"{index}: {r[0]} - {r[1]}\n"
            index += 1

        return output

    def id_is_valid(self, id: int):
        logger.debug(f"checking id {id}")
        index = 0
        for r in self.valid_ranges:
            if r[0] > id:
                logger.debug(f"id is below range {index}: {r}")
                return False
            if r[0] <= id and r[1] >= id:
                # logger.debug(f"id is in range {index}: {r}")
                return True
            index += 1
        raise Exception("just curious if this can happen")

    @staticmethod
    def from_ranges(ranges: list[tuple[int, int]]):
        sorted_ranges = sorted(ranges, key=lambda item: item[0])

        logger.debug("sorted ranges:")
        for r in sorted_ranges:
            logger.debug(f"{r[0]} - {r[1]}")

        final_ranges: list[tuple[int, int]] = []

        logger.debug("merging ranges ...")
        current_range = sorted_ranges[0]
        logger.debug(f"current_range = {current_range[0]} - {current_range[1]}")
        for r in sorted_ranges:
            logger.debug(f"r = {r[0]} - {r[1]}")
            if r[0] <= current_range[1] and r[1] >= current_range[1]:
                logger.debug("extending current_range ...")
                current_range = (current_range[0], r[1])
            else:
                logger.debug("starting new range ...\n")
                final_ranges.append(current_range)
                current_range = r
            logger.debug(f"current_range = {current_range[0]} - {current_range[1]}")

        final_ranges.append(current_range)

        output = database(final_ranges)
        return output

    @staticmethod
    def from_file(file: TextIOBase):
        ranges: list[tuple[int, int]] = []

        for line in file:
            line = line.strip()
            if line == "":
                break
            parts = line.split("-")
            ranges.append((int(parts[0]), int(parts[1])))

        return database.from_ranges(ranges)

    def validate_ids_from_file(self, file: TextIOBase):
        logger.debug("validating ids from file")
        logger.debug(self)
        total = 0

        for line in file:
            line = line.strip()
            if line == "":
                break
            id = int(line)
            if self.id_is_valid(id):
                total += 1

        return total

    def count_valid_ids(self):
        total = 0
        for r in self.valid_ranges:
            ids_in_range = r[1] - r[0] + 1
            total += ids_in_range 
            logger.debug(f"{r[1]} - {r[0]} + 1 = {ids_in_range}")

        return total


def test():
    input = """3-5
10-14
16-20
12-18

1
5
8
11
17
32"""

    string_as_file = StringIO(input)
    db = database.from_file(string_as_file)

    total = db.validate_ids_from_file(string_as_file)

    print(f"valid ids: {total}")


def test_custom():
    input = """5-15
20-30
3-10

1
5
10
15
18
20
25
30
35
        """

    string_as_file = StringIO(input)
    db = database.from_file(string_as_file)

    total = db.validate_ids_from_file(string_as_file)

    print(f"valid ids: {total}")


def main():
    length = len(argv)
    if length != 2:
        raise Exception(f"Expected 1 argument, got {length - 1}.")

    with open(argv[1]) as f:
        db = database.from_file(f)
        total = db.validate_ids_from_file(f)

    print(f"valid ids from file: {total}")
    print(f"valid ids in ranges: {db.count_valid_ids()}")


if __name__ == "__main__":
    logging.basicConfig(level=logging.DEBUG)

    main()
