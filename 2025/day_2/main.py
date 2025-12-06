from sys import argv


def id_is_invalid(id: int):
    id_as_string = str(id)

    id_length = len(id_as_string)

    pattern_length = 0
    while (pattern_length + 1) * 2 <= id_length:
        pattern_length += 1

        if id_length % pattern_length != 0:
            continue

        pattern = id_as_string[0:pattern_length]

        pattern_occurence = id_length // pattern_length

        if id_as_string == pattern * pattern_occurence:
            return True

    return False


def get_ranges_from_input(input_string: str):
    ranges = input_string.split(",")
    output: list[tuple[int, int]] = []
    for r in ranges:
        parts = r.split("-")
        if len(parts) != 2:
            raise ValueError("malformed range")
        output.append((int(parts[0]), int(parts[1])))

    return output


def get_total_from_range(r: tuple[int, int]):
    total = 0
    for id in range(r[0], r[1] + 1):
        if id_is_invalid(id):
            # print(f"{id} is invalid")
            total += id

    return total


def test():
    input = "11-22,95-115,998-1012,1188511880-1188511890,222220-222224,1698522-1698528,446443-446449,38593856-38593862,565653-565659,824824821-824824827,2121212118-2121212124"
    ranges = get_ranges_from_input(input)

    total = 0
    for r in ranges:
        total += get_total_from_range(r)

    print(f"the total is {total}")
    if total != 4174379265:
        raise Exception("total doesn't match expected total")


def get_total_ids_from_input_file():
    if len(argv) != 2:
        raise ValueError(f"Expected 1 argument, got {len(argv)}.")

    file_path = argv[1]
    with open(file_path) as f:
        file_content = f.read()

    ranges = get_ranges_from_input(file_content)
    total = 0
    for r in ranges:
        total += get_total_from_range(r)

    print(f"total is {total}")


if __name__ == "__main__":
    get_total_ids_from_input_file()
