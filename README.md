### Simplex Algorithm
A linear programming solver (e.g., cost minimization) built from scratch using the Simplex method.

---

#### Overview:
Take text-based inputs in the form:
```
let x1;
let x2;
max 3 * x1 + 4 * x2;
s.t. x1 + x2 <= 5;
x1 >= 0;
x2 >= 0;
```
and receive optimal solutions (or determine infeasibility/unboundedness).
<p align="left">
  <img src="https://drive.google.com/uc?export=view&id=1K2ajupc39jTZ8I4frbeb57J4DaW_87wz" alt="Simple Mail Preview" width="762" style="border:10px solid #ddd;" />
</p>

---

#### How it works:
```mermaid
graph LR
    subgraph Frontend
        A("Input<br>(TypeScript, Next.js, React)")
    end

    subgraph Go_Pipeline["Go Pipeline"]
        B("Backend API<br>(Go net/http)")
        C("Lexer / Parser / Simplifier<br>+ Expression Optimizations<br>(Go)")
    end

    subgraph Simplex_Algorithm["C++ Algorithm (Math)"]
        D("Standard Equality Form")
        E("Canonical Form")
        F("2-Phase Simplex")
        G("Custom Matrix Class")
    end

    A --> B --> C --> D --> E --> F
```
