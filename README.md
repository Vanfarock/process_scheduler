# Simulation of an OS process scheduler

## Diagram of the project
![Process states graph](https://github.com/Vanfarock/process_scheduler/blob/master/processes_states_graph.jpg)

## Result of execution
![Visual result of the execution](https://github.com/Vanfarock/process_scheduler/blob/master/visual_result.png)

## Features
- [X] Moving processes between states, following the above diagram
- [X] Spawn x threads for each state, where x is the number of cores in the simulated CPU
- [X] Identify IO-boundness of each process
- [X] Prioritize CPU-bound processes over IO-bound
- [X] Display history of execution in detailed/visual way
