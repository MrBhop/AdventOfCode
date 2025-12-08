from io import StringIO, TextIOBase
from sys import argv
from typing import override


class database:
    def __init__(self, valid_ranges: list[tuple[int, int]] | None = None) -> None:
        self.valid_ranges: list[tuple[int, int]] = valid_ranges or []

    @override
    def __repr__(self) -> str:
        output = "database:\n"

        index = 0
        for range in self.valid_ranges:
            output += f"{index}: {range[0]} - {range[1]}\n"
            index += 1

        return output

    def id_is_valid(self, id: int):
        print(f"checking id {id}")
        index = 0
        for range in self.valid_ranges:
            if range[0] > id:
                print(f"id is below range {index}: {range}")
                return False
            if range[0] <= id and range[1] >= id:
                print(f"id is in range {index}: {range}")
                return True
            index += 1
        raise Exception("just curious if this can happen")

    @staticmethod
    def from_ranges(ranges: list[tuple[int, int]]):
        sorted_ranges = sorted(ranges, key=lambda item: item[0])

        final_ranges: list[tuple[int, int]] = []

        current_range = sorted_ranges[0]
        for range in sorted_ranges:
            if range[0] <= current_range[1]:
                current_range = (current_range[0], range[1])
            else:
                final_ranges.append(current_range)
                current_range = range

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
        print("validating ids from file")
        print(self)
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

    print(f"valid ids: {total}")


if __name__ == "__main__":
    length = len(argv)
    if length != 2:
        raise Exception(f"Expected 1 argument, got {length - 1}.")

    with open(argv[1]) as f:
        ranges: list[tuple[int, int]] = []

        for line in f:
            line = line.strip()

            if line == "":
                break
            parts = line.split("-")
            ranges.append((int(parts[0]), int(parts[1])))

        total = 0
        for line in f:
            id = int(line.strip())

            for range in ranges:
                if range[0] <= id and range[1] >= id:
                    total += 1
                    break

    print(f"valid ids: {total}")
