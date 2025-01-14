eBPF allows application code to run in pre-compiled form (bytecode) at the level of the Linux kernel. Running as bytecode in the Linux kernel significantly enhances app performance, as the interpreter gets cut out of the equation and the app code can directly access system resources without having to traverse the OSI stack. Applications can directly listen to system-level events such as…

…network traffic: knowing the type, number, size, and origin of incoming and outgoing packets is key for optimizing network configurations, detecting security threats, and planning future upgrades.

…filesystem I/O: file operations such as open, close, read, and write can be correlated with other events, such as network events, to gain insights into app usage patterns and how they impact the rest of the application.

…system calls: watching system calls made by applications to learn how they interact with different parts of the application stack, such as the Kubernetes scheduler, Docker containers, or VMs.

…process lifecycle: understand application behavior through following the creation, execution, and termination of their processes. This becomes especially important when operating distributed apps that depend on numerous microservices working together.

…resource utilization: Memory allocation, CPU scheduling events, and disk I/O are all observable with eBPF. This real-time data is crucial for performance tuning and capacity planning.

…security monitoring: eBPF is capable of detecting changes to critical files, monitoring for common exploit signatures, and tracking user-level login events, making it an ideal tool for intrusion detection systems.

…kernel tracing: eBPF can listen to kernel tracepoints, which are static hooks in the kernel, to gather metrics or logs for debugging and performance analysis without impacting system performance significantly.

…hardware events: Through Performance Monitoring Counters (PMCs), eBPF can collect data on hardware-level events like CPU cache misses or memory paging aiding in low-level performance analysis.

…custom metrics: Beyond standard system events, eBPF allows for creating custom metrics based on specific needs. Users can write eBPF programs that define exactly what data to collect, how to aggregate it, and when to report it.

…scheduling and threading: Observations on how the kernel scheduler operates, how threads are created and destroyed, and how context switches occur are possible with eBPF. This is crucial for understanding the performance of concurrent applications.

The ability to directly listen to the Linux kernel’s event stream significantly reduces the need for instrumentation and therefore minimizes the risk originating from unmonitored systems. Correlating all of these event streams is the foundation for the automatic creation and continuous updating of a comprehensive dependency map for the entire organization.

