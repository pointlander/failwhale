# FailWhale

The FailWhale algorithm determines the probability that a request should be blocked and not forwarded based on the success, failure, or being block of past requests. The overall architecture of the algorithm is an artificial neuron with time differences as the input. The weights are based on whether past requests were a success (1), failure (-1), or blocked (0). A circular buffer of past requests is maintained. Each entry of the circular buffer has a time stamp and a weight. The time differences are calculated first:

```
Di = T - Si
```

T is the current time, and S is the time stamp. The next step is to normalize the time differences, D, using the [SoftMax](https://en.wikipedia.org/wiki/Softmax_function) function:

```
NormalizedDi = e^(-Di/Max(Di)) / Sum(e^(-Di/Max(Di)))
```

Finally, each NormalizedDi is multiplied by a weight, summed, and fed into a [sigmoid](https://en.wikipedia.org/wiki/Sigmoid_function):

```
Probability = 1 / (1 + e^(8 * Sum(NormalizedDi * Wi)))
```

Whether or not an incoming request is blocked is determined using the Probability (and a random number generator). If the request is blocked, then an entry is added to the circular buffer with a weight of 0 and the current time for the time stamp. If the request is not blocked, then an entry is added to the circular buffer with a weight of 1 or -1 (depending on if the forwarded request is successful or a failure) and the current time for the time stamp.

Below is the probability output (the probability that a request is blocked) for a scenario:

![scenario](probability_vs_time.png?raw=true)

Firstly, all requests are a success, then all failures, then a mixture of success and failures, and finally all successes.
