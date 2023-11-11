#include <stdio.h>
#include <stdlib.h>
#include <mpi.h>
#include <stdbool.h>

#define ITERATIONS 500

int main(void)
{
    // i is the rank of the process to be passed during iterations
    // n represents the number of processes to be spawned
    int n, i;

    MPI_Init(NULL, NULL);

    MPI_Comm_size(MPI_COMM_WORLD, &n);
    MPI_Comm_rank(MPI_COMM_WORLD, &i);

    // when the process becomes neutral, we will override isNeutral to true
    bool isNeutral = false;

    int index;
    for (index = 0; index < ITERATIONS; index++)
    {
        // for even iterations
        if (index % 2 == 0)
        {
            int highestNumber = -1;

            // instead of each process sending its rank to one another and then performing MAX(i),
            // use collective communication to extract i's from each process, compute MPI_MAX, broadcast it to all processes as highestNumber
            // Allreduce function allows us to achieve above functionality
            MPI_Allreduce(&i, &highestNumber, 1, MPI_INT, MPI_MAX, MPI_COMM_WORLD);

            // if process is NOT already neutral AND highest number is same as the process rank
            if (!isNeutral && highestNumber == i)
            {
                // mark process as neutral
                isNeutral = true;
                printf("Even Iteration %d: Marking process(%d) as neutral\n", index, i);
            }
        }
        else
        {
            // determine neighbor processes
            int nextProcess = i + 1, prevProcess = i - 1;
            if (n == 1) {
                break; // there's no neighbor to send the data to
            }
            else if (i == 0)
            {
                prevProcess = n - 1;
            }
            else if (i == n - 1)
            {
                nextProcess = 0;
            }

            int recvBuffer;
            // every process will send 'i' value to next process (Tag=0) and receive the 'i-1' value from previous process (Tag=0)
            MPI_Sendrecv(&i, 1, MPI_INT, nextProcess, 0, &recvBuffer, 1, MPI_INT, prevProcess, 0, MPI_COMM_WORLD, MPI_STATUS_IGNORE);

            // if the received value is 0, mark the process as neutral
            if (!isNeutral && recvBuffer == 0)
            {
                isNeutral = true;
                printf("Odd Iteration %d: Marking process(%d) as neutral\n", index, i);
            }

            int sendBuffer = recvBuffer;
            // if process is NOT neutral, subtract 1 from the value received
            if (!isNeutral)
            {
                sendBuffer--;
            }

            // send this subtracted/not subtracted (in case of neutral) value to the next process (Tag=1)
            // and receive the similarly processed values from the previous process (Tag=1)
            MPI_Sendrecv(&sendBuffer, 1, MPI_INT, nextProcess, 1, &recvBuffer, 1, MPI_INT, prevProcess, 1, MPI_COMM_WORLD, MPI_STATUS_IGNORE);

            // if the received value is 0, mark the process as neutral
            if (!isNeutral && recvBuffer == 0)
            {
                isNeutral = true;
                printf("Odd Iteration %d: Marking process(%d) as neutral\n", index, i);
            }
        }
    }

    MPI_Barrier(MPI_COMM_WORLD);

    if (isNeutral)
    {
        printf("Process %d is neutral\n", i);
    }

    MPI_Finalize();

    return 0;
}