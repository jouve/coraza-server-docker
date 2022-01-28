import asyncio
import logging
import sys
from argparse import ArgumentParser
from asyncio import Task

from httpx import AsyncClient


async def amain(hostport: str, parralel: int, count: int):
    async with AsyncClient() as session:
        total = int(0)
        tasks = set[Task]()
        while True:
            todo = min(count - total - len(tasks), parralel - len(tasks))
            tasks.update(
                asyncio.create_task(session.get(f'http://{hostport}/'))
                for _ in range(todo)
            )
            await asyncio.sleep(0.1)
            done, tasks = await asyncio.wait(tasks, return_when=asyncio.FIRST_COMPLETED)
            total += sum(1 if result.exception() else 1 for result in done)
            logging.info('running=%s/%s done=%d/%d', len(tasks), parralel, total, count)
            if total >= count:
                break


def main():
    parser = ArgumentParser()
    parser.add_argument('-V', '--version', action='version', version='0.1.0')
    parser.add_argument('-v', '--verbose', action='count', default=0)
    parser.add_argument('hostport')
    parser.add_argument('parallel', type=int)
    parser.add_argument('count', type=int)
    args = parser.parse_args()

    logging.basicConfig(level=[logging.WARNING, logging.INFO, logging.DEBUG][min(args.verbose, 2)])

    return asyncio.run(amain(args.hostport, args.parallel, args.count))


if __name__ == '__main__':
    sys.exit(main())
