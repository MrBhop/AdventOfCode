from dataclasses import dataclass
from io import StringIO, TextIOBase
import logging
from sys import argv
from typing import override


logger = logging.getLogger(__name__)


@dataclass
class IdRange:
    lower_bound: int
    upper_bound: int

    def count_valid_ids(self):
        return self.upper_bound - self.lower_bound + 1

    @staticmethod
    def from_string(input: str):
        parts = input.split("-")
        return IdRange(int(parts[0]), int(parts[1]))

    @override
    def __repr__(self) -> str:
        return f"IdRange({self.lower_bound}-{self.upper_bound})"


def merge_ranges(id_ranges: list[IdRange]):
    """Returns a new list of sorted and merged ranges, covering the same numbers as the input ranges."""

    logger.debug("merging ranges")

    # sort the input by lower_bound.
    sorted_ranges = sorted(id_ranges, key=lambda item: item.lower_bound)

    final_ranges: list[IdRange] = []
    current_range = sorted_ranges[0]

    for r in sorted_ranges:
        logger.debug("current_range:")
        logger.debug("%s", current_range)
        logger.debug("r:")
        logger.debug("%s", r)
        if r.lower_bound <= current_range.upper_bound + 1:
            # r overlaps with current_range. Merge 'em.
            current_range.upper_bound = max(current_range.upper_bound, r.upper_bound)
            logger.debug("merged range.")
        else:
            # r does not overlap with current_range. Add current_range to the output and use r for the next loop.
            final_ranges.append(current_range)
            current_range = r
            logger.debug("starting new range.\n")
    if current_range.upper_bound != final_ranges[-1].upper_bound:
        # if the leftover current_range is not already in the output, add it as well.
        final_ranges.append(current_range)
    return final_ranges


@dataclass
class Database:
    id_ranges: list[IdRange]

    @override
    def __repr__(self) -> str:
        return f"Database({',\n'.join([repr(item) for item in self.id_ranges])})"

    def id_is_valid(self, id: int):
        for r in self.id_ranges:
            if r.lower_bound <= id and r.upper_bound >= id:
                logger.debug("id %d is in range %s", id, r)
                return True
        logger.debug("id %d is not in any range.", id)
        return False

    def count_valid_ids(self):
        total = 0
        for r in self.id_ranges:
            total += r.count_valid_ids()
        return total

    @staticmethod
    def from_ranges(id_ranges: list[IdRange]):
        merged_ranges = merge_ranges(id_ranges)

        return Database(merged_ranges)

    @staticmethod
    def from_file(file: TextIOBase):
        ranges: list[IdRange] = []

        for line in file:
            line = line.strip()

            if line == "":
                break
            ranges.append(IdRange.from_string(line))
        return Database.from_ranges(ranges)

    def validate_ids_from_file(self, file: TextIOBase):
        total = 0

        for line in file:
            line = line.strip()

            if line == "":
                break
            id = int(line)
            if self.id_is_valid(id):
                total += 1
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
    db = Database.from_file(string_as_file)

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
    db = Database.from_file(string_as_file)

    total = db.validate_ids_from_file(string_as_file)

    print(f"valid ids: {total}")


def main():
    length = len(argv)
    if length != 2:
        raise Exception(f"Expected 1 argument, got {length - 1}.")

    with open(argv[1]) as f:
        db = Database.from_file(f)
        total = db.validate_ids_from_file(f)

    print(f"valid ids from file: {total}")
    print(f"valid ids in ranges: {db.count_valid_ids()}")


if __name__ == "__main__":
    logging.basicConfig(level=logging.DEBUG)

    main()
