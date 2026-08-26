<p align="center">
  <img src="assets/atlas-banner.png" alt="Atlas Infrastructure Intelligence" width="100%">
</p>

# Atlas

**Infrastructure Intelligence by Iron Signal Systems**

> **Truth first. Built for operators.**

Atlas is a read-focused infrastructure intelligence platform designed to reconstruct how a network is actually built, connected, controlled, and changing.

It collects approved information from infrastructure systems, preserves where that information came from, normalizes and correlates the facts, builds an infrastructure model, and uses that model to answer operational questions and generate current network diagrams.

Atlas does **not** ingest network diagrams to determine network topology.

Atlas generates diagrams from the same infrastructure model used to answer questions.

```text
Infrastructure Sources
        ↓
      Collect
        ↓
Preserve Source Identity
        ↓
     Normalize
        ↓
     Correlate
        ↓
       Model
      /     \
     /       \
 Answers    Diagrams
```

> **Do not make the operator browse the network. Reconstruct the network and answer the question.**

---

## Development Status

**Pre-alpha. Not production ready.**

This repository is a clean restart of Atlas.

The previous implementation accumulated more architecture, project structure, abstractions, tooling, and process than the product required at its current stage.

That implementation is retained at:

[Iron-Signal-Systems/old-atlas](https://github.com/Iron-Signal-Systems/old-atlas)

The old repository is historical reference material.

Useful code, tests, parsers, SQL, validation cases, and design ideas may be brought forward selectively.

Nothing is migrated simply because it existed before.

The current Atlas implementation follows a simpler rule:

> **Build the simplest correct thing that answers a real infrastructure question.**

Complexity is allowed when the infrastructure problem requires it.

Complexity must earn its place.

A controlled pilot is currently targeted approximately **24–36 months** from this restart, subject to engineering and validation results.

---

# What Atlas Is

Atlas is an **infrastructure intelligence platform**.

It is intended to understand relationships across network infrastructure so operators do not have to repeatedly reconstruct those relationships by hand.

Atlas may eventually understand information such as:

* devices;
* interfaces;
* VLANs;
* trunks;
* port channels;
* neighbors;
* MAC and ARP observations;
* IP addresses and subnets;
* routing domains;
* routes and next hops;
* firewall policies;
* ACLs;
* NAT;
* VPNs;
* SD-WAN;
* wireless infrastructure;
* sites and physical placement;
* dependencies;
* historical changes;
* identity relationships from specialist systems; and
* other infrastructure facts required to answer useful questions.

These capabilities are added when real Atlas questions require them.

Atlas does not need to understand every command, feature, protocol, or vendor capability simply because the information exists.

---

# The Problem Atlas Solves

Experienced network administrators and engineers can usually answer difficult infrastructure questions.

The problem is the amount of work required to reach the answer.

A single troubleshooting question may require:

* connecting to multiple switches or routers;
* checking interface configuration;
* tracing VLANs;
* reviewing spanning tree;
* checking MAC and ARP tables;
* examining routing tables;
* inspecting firewall policies;
* checking monitoring systems;
* comparing old documentation;
* looking through tickets or change records; and
* relying on the memory of experienced staff.

Atlas is intended to perform that correlation once and make the result reusable.

Questions Atlas should eventually answer include:

* Where is this IP address?
* Which switch and interface is it connected to?
* Which VLAN contains it?
* Which subnet and gateway apply?
* Which route would traffic use?
* Why was that route selected?
* Can source A reach destination B on a specific protocol or port?
* Which firewall policy or ACL affects that path?
* Is NAT involved?
* Is the return path valid?
* What depends on this switch?
* What depends on this interface?
* What depends on this circuit?
* What depends on this VLAN?
* What happens if this device or site becomes unavailable?
* Does a supposed redundant path actually support the required service?
* What changed before an outage began?
* Which source observations support an answer?
* Which parts of the answer are incomplete, stale, inferred, conflicting, or unknown?

> **What takes an experienced engineer hours to reconstruct manually should, when Atlas has the required information, be answerable in seconds or minutes.**

---

# Atlas Generates Network Diagrams

Network diagrams are an **output of Atlas**, not an authoritative topology input.

Atlas reconstructs topology from approved infrastructure observations and generates diagrams from the resulting model.

```text
Cisco
FortiGate
Other approved infrastructure sources
        ↓
      Atlas
        ↓
Infrastructure Model
        ↓
Generated Network Diagram
```

Atlas-generated diagrams may eventually show information such as:

* sites;
* network devices;
* device-to-device links;
* switch relationships;
* trunks;
* port channels;
* VLAN propagation;
* routing relationships;
* WAN connections;
* firewall boundaries;
* wireless infrastructure;
* circuit dependencies;
* endpoints; and
* other relationships represented by the Atlas model.

Generated diagrams should reflect the same knowledge Atlas uses to answer operational questions.

A diagram must not become a second independently maintained source of network truth.

If Atlas knowledge changes, the generated diagram should change with it.

If Atlas does not know something, the diagram should not invent it.

---

# Diagrams Are Not Topology Sources

Atlas does not determine network truth from an existing Draw.io, Visio, PDF, image, or manually maintained network diagram.

Human-created diagrams may be useful to people for planning or reference, but they are not authoritative Atlas topology inputs.

Atlas should not reason:

```text
diagram says switch A connects to switch B
        ↓
therefore the connection exists
```

Instead Atlas should establish the relationship from approved infrastructure information.

For example:

```text
Cisco CDP / LLDP
interface configuration
port-channel state
VLAN configuration
other relevant observations
        ↓
Atlas correlation
        ↓
switch A ↔ switch B
```

Diagram formats such as Draw.io may eventually be supported as **publication or export formats**.

---

# Practical Example

Assume a fiber failure takes down a camera network at a building that was believed to have a redundant path.

The obvious fact may be:

```text
fiber failed
```

The useful question is:

> **Why did this fiber failure take down the camera VLAN if the site was supposed to be redundant?**

The actual problem may have been introduced months earlier:

* a switch was replaced;
* the alternate topology changed;
* the camera VLAN was never extended across the alternate path;
* normal traffic continued using the primary link;
* the missing dependency remained hidden; and
* the failure exposed the problem.

Atlas should eventually be able to explain:

* which switches participate;
* which links connect them;
* which VLANs cross those links;
* which path normally carries the traffic;
* whether the alternate path supports the affected VLAN;
* which routing and firewall controls participate;
* when the relationship changed;
* what else depends on the same path; and
* which source observations support the conclusion.

The objective is not simply to create a prettier network diagram.

The objective is to make network relationships understandable.

---

# Core Operating Model

Atlas follows a simple conceptual flow:

```text
COLLECT
   ↓
PRESERVE SOURCE IDENTITY
   ↓
PARSE
   ↓
NORMALIZE
   ↓
CORRELATE
   ↓
MODEL
   ↓
REASON
   ↓
ANSWER / MAP / COMPARE / EXPLAIN
```

This is a **data-flow model**, not a requirement that every stage become a separate service, framework, interface, package, or executable.

Several stages may initially exist in the same small implementation.

Architecture should be extracted when working code demonstrates that the separation is useful.

---

# Current Implementation Focus

The current implementation focus is **Cisco source parsing**.

The Cisco parser is being rebuilt carefully rather than copied wholesale from the previous Atlas implementation.

The objective is not to parse every Cisco command.

The objective is to faithfully collect the Cisco facts Atlas needs to answer useful infrastructure questions.

Early Cisco work may include information such as:

* device identity;
* platform information;
* interfaces;
* interface addresses;
* interface state;
* VLANs;
* trunks;
* port channels;
* CDP;
* LLDP;
* MAC observations;
* ARP observations;
* spanning tree;
* routes;
* VRFs;
* ACL-related information; and
* other facts required by current Atlas questions.

Features are added because Atlas needs the information.

They are not added simply because Cisco exposes the command.

---

# Parser Requirements

Source parsers are a critical trust boundary.

A parser should preserve what the source actually reported before Atlas attempts higher-level interpretation.

The initial rule is:

> **Parse faithfully first. Correlate later. Derive last.**

A parser should not quietly turn assumptions into source facts.

For example:

```text
Source:
Interface Gi1/0/24 allows VLANs 10,20,30

Parser:
records VLANs 10,20,30 as configured on the trunk
```

The parser should not automatically claim:

```text
VLAN 20 is currently carrying traffic
```

unless a source observation supports that conclusion.

Parser behavior should be:

* bounded;
* deterministic where practical;
* testable;
* explicit about unsupported input;
* explicit about malformed input;
* conservative when information is incomplete; and
* resistant to hostile or unexpectedly large source content.

Unknown input should not automatically crash the entire collection.

Unsupported input should be identifiable so support can be added deliberately.

---

# Source-Backed Information

Atlas should preserve enough information to establish where important facts came from.

Depending on the source and record type, this may include:

* source system;
* source device;
* source command or collection operation;
* collection timestamp;
* parser identity;
* parser version;
* source identifiers;
* source values;
* cryptographic source identity where useful;
* direct observations;
* derived values;
* assumptions;
* conflicting information; and
* source freshness.

Atlas should distinguish between different kinds of knowledge.

Initial state vocabulary may include:

* `CONFIGURED`
* `OBSERVED`
* `CALCULATED`
* `INFERRED`
* `UNKNOWN`
* `CONFLICTING`

A configuration proves configuration.

It does not automatically prove current operational state.

A route being configured does not prove that it is selected.

A VLAN being permitted on a trunk does not prove that traffic is currently using it.

A Layer-3 path existing does not prove that a firewall permits the service.

A missing observation does not automatically prove absence.

> **If Atlas cannot establish something, Atlas should say that it cannot establish it.**

---

# Infrastructure Model

Atlas should model infrastructure relationships rather than reproduce vendor navigation structures.

Vendor-specific information remains available for traceability, but vendor syntax should not become the canonical Atlas model.

Potential model areas include the following.

## Device and placement

* organizations;
* sites;
* buildings;
* rooms;
* racks;
* devices;
* chassis;
* stacks;
* logical devices; and
* components.

## Layer 2

* interfaces;
* VLANs;
* trunks;
* native VLANs;
* allowed VLANs;
* port channels;
* neighbors;
* MAC observations;
* spanning-tree relationships; and
* endpoint attachment.

## Layer 3

* IP addresses;
* prefixes;
* subnets;
* gateways;
* VRFs;
* VDOMs;
* routing domains;
* routes; and
* next hops.

## Network controls

* firewall zones;
* firewall policies;
* ACLs;
* address objects;
* service objects;
* NAT;
* VIPs;
* VPNs;
* SD-WAN; and
* management boundaries.

These model areas are not instructions to create all corresponding code immediately.

The model grows as working Atlas capabilities require it.

---

# Read-Only by Design

Atlas is designed around **read-only collection and least privilege**.

The normal Atlas operating direction is:

```text
read
  ↓
collect
  ↓
parse
  ↓
normalize
  ↓
correlate
  ↓
reason
  ↓
answer
```

Atlas does not require administrative control over infrastructure simply to understand it.

Collectors should receive only the authority necessary to retrieve approved source information.

Atlas should not silently modify infrastructure.

Any future write, provisioning, or remediation capability would require a separately engineered authority boundary.

That capability would need explicit authorization, attribution, bounded scope, validation, and recovery behavior.

Read-only operation is not a missing feature.

It is part of the Atlas trust model.

---

# Specialist Tools Remain Specialists

Atlas is not intended to replace systems that already perform specialized functions well.

Examples include:

* Zabbix and similar monitoring systems;
* Graylog and similar logging platforms;
* Security Onion and other network-security platforms;
* CrowdStrike and other endpoint-security products;
* BloodHound for identity and privilege relationship analysis;
* Cisco infrastructure platforms;
* Fortinet infrastructure platforms; and
* other specialist operational systems.

Atlas may consume approved information from specialist systems where that information contributes to infrastructure understanding.

Atlas adds correlation across those domains.

> **Specialist tools remain specialists. Atlas fills the gaps between them.**

---

# Future Source Work

Cisco is the current source implementation.

Other sources may be added when Atlas requires them.

Potential future sources include:

## FortiGate

Useful information may eventually include:

* device identity;
* interfaces;
* VLANs;
* zones;
* VDOMs;
* firewall policies;
* objects;
* NAT;
* VIPs;
* routing;
* VPNs;
* SD-WAN;
* HA state; and
* management-plane controls.

## BloodHound

BloodHound may eventually provide identity and privilege context.

Atlas should preserve BloodHound semantics rather than reinterpret them during source collection.

```text
BloodHound says X
        ↓
Atlas records X
        ↓
later correlation
```

Identity relationships and network reachability are different facts.

Atlas should not collapse them into one unsupported conclusion.

## Other Sources

Additional infrastructure, monitoring, security, wireless, asset, or operational sources may be supported when they provide facts needed to answer real Atlas questions.

There is no requirement to create a generic plugin framework before those sources exist.

---

# Dependency Analysis

Network topology is only part of the problem.

Atlas should eventually understand dependencies.

Useful questions may include:

* What depends on this switch?
* What depends on this interface?
* What depends on this firewall?
* What depends on this route?
* What depends on this circuit?
* What depends on this VLAN?
* What stops working if this site disappears?
* Which services share a single point of failure?
* Does an alternate path actually support everything the service requires?

Dependencies should be derived from Atlas knowledge rather than maintained manually as a second truth source.

---

# Historical State

Atlas is intended to preserve historical infrastructure knowledge.

New observations should not silently erase older observations.

This should eventually allow questions such as:

* What did the network look like before this change?
* When did this route appear?
* When did this VLAN stop crossing this trunk?
* Was this interface active during the outage?
* What changed between two collections?
* Did the post-change network match the expected result?

History should remain understandable and attributable.

> **New knowledge should not silently rewrite old knowledge.**

Implementation of historical storage should remain as simple as possible while preserving this property.

---

# PostgreSQL Direction

PostgreSQL is the planned persistent Atlas data store.

The database design should support:

* source-backed records;
* infrastructure relationships;
* historical observations;
* correlation;
* query;
* generated topology;
* source identity; and
* integrity validation where required.

The schema should grow from working data requirements.

Atlas should not attempt to design the final database before the required records and relationships have been proven through implementation.

Database complexity must earn its place just like application complexity.

---

# Security Principles

Atlas should be built with a small and understandable attack surface.

Initial security priorities include:

* read-only collection;
* least privilege;
* explicitly scoped collection;
* untrusted-input handling;
* bounded parsers;
* resource limits;
* timeouts where appropriate;
* no secrets committed to Git;
* no unrestricted production source data committed to Git;
* conservative parsing;
* explicit uncertainty;
* clear trust boundaries;
* tested failure behavior;
* dependency review;
* reproducible behavior where practical; and
* straightforward code that can be reviewed.

Security controls should solve real risks.

Security architecture should not become unnecessary framework complexity.

> **Simple code that can be understood and reviewed is itself a security property.**

---

# Development Principles

## Build the simplest correct thing

> **Build the simplest correct thing that answers a real Atlas question.**

Do not build infrastructure for hypothetical future requirements.

---

## Complexity must earn its place

Infrastructure is complicated.

The implementation does not need to make that complication worse.

A complex problem should be broken into understandable steps.

---

## Source truth before interpretation

Preserve source facts before deriving conclusions from them.

---

## Parse faithfully first

> **Parse faithfully first. Correlate later. Derive last.**

---

## Unknown is valid

Atlas should prefer:

```text
UNKNOWN
```

over a confident answer that cannot be supported.

---

## Tests over ceremony

Engineering confidence should come primarily from:

* unit tests;
* integration tests where useful;
* deterministic fixtures;
* hostile-input testing;
* fuzzing where valuable;
* malformed-input testing;
* resource-bound testing;
* failure testing;
* recovery testing;
* security review; and
* validation against real authorized infrastructure.

Documentation should exist when it helps build, operate, secure, support, test, or explain Atlas.

Process does not substitute for correctness.

---

## No speculative frameworks

Do not create a generalized framework because it might be useful someday.

When multiple working implementations demonstrate a real common contract, extract the common behavior.

Before then, duplication may be cheaper and safer than the wrong abstraction.

---

## Prefer obvious code

Atlas code should favor:

* small functions;
* clear names;
* explicit control flow;
* understandable data structures;
* minimal hidden behavior;
* limited package coupling; and
* straightforward error handling.

Clever code should require justification.

---

## Keep responsibilities separate

Collection, parsing, correlation, storage, and presentation are different responsibilities.

They do not necessarily require separate services or frameworks.

Separation should exist where it improves understanding, testing, security, or maintainability.

---

## Git already remembers

Obsolete architecture, abandoned roadmaps, and historical process structures do not need to remain in the active project simply to preserve history.

Git and `old-atlas` already provide that history.

---

# Code Migration From old-atlas

Code should be migrated by **capability**, not by directory.

Do not reason:

```text
old-atlas has this package
therefore new Atlas needs this package
```

Instead ask:

```text
new Atlas requires capability X
        ↓
does old-atlas contain a simple,
correct, understood implementation?
        ↓
yes
        ↓
bring forward the smallest useful portion
plus its tests
```

Useful old implementation work should not be rewritten solely because it is old.

It should be:

1. understood;
2. reviewed;
3. simplified where necessary;
4. tested;
5. migrated in the smallest useful unit; and
6. used only if the new Atlas currently needs it.

Old architecture does not automatically become new architecture.

---

# Repository Structure

The repository should remain small while Atlas is being rebuilt.

Directories are created when implementation requires them.

An early structure may be as small as:

```text
.
├── cmd/
│   └── atlas/
├── internal/
│   ├── cisco/
│   └── records/
├── assets/
├── README.md
├── LICENSE
├── SECURITY.md
└── go.mod
```

This is illustrative, not mandatory.

As real capabilities appear, additional areas may be introduced.

For example:

```text
internal/
├── cisco/
├── fortigate/
├── records/
├── model/
└── postgres/
```

A generic source or module framework should not be created until multiple implementations demonstrate that it is needed.

The repository structure should follow the implementation.

The implementation should not be forced to follow a speculative repository structure.

---

# Initial Development Direction

Atlas development should proceed through working capability rather than large vendor-numbered phases.

A practical direction is:

```text
source
   ↓
faithful parsing
   ↓
source-backed records
   ↓
small infrastructure model
   ↓
correlation
   ↓
answer a useful question
   ↓
generate a view of the same model
```

The current Cisco work should prove this path first.

For example, an early Atlas question may be:

> **Where is this IP address?**

Answering that may require Atlas to understand:

* device identity;
* interfaces;
* VLANs;
* MAC observations;
* ARP observations;
* IP addressing;
* subnets;
* neighbor relationships; and
* relevant source freshness.

Those requirements then drive implementation.

The next question drives the next capability.

---

# Product Test

Every significant Atlas capability should face the same test:

> **Does this save a qualified operator from manually correlating infrastructure information to reach the same conclusion?**

Progress is not measured by:

* parser count;
* package count;
* interface count;
* service count;
* dashboard count;
* document count;
* roadmap phase count; or
* raw database record count.

Progress is measured by:

* useful questions answered;
* correctness;
* visible uncertainty;
* time saved;
* dependencies understood;
* changes explained;
* operational effort reduced; and
* operator trust.

---

# Pilot Direction

The first controlled Atlas pilot should validate the system against a real authorized environment.

The pilot should determine whether:

* Atlas answers are correct;
* missing information is represented truthfully;
* source information remains traceable;
* troubleshooting time is reduced;
* network diagrams remain useful and current;
* hidden dependencies become visible;
* collection is safe and low-impact;
* another qualified administrator can understand the environment without relying on one person's memory;
* the system remains maintainable; and
* failures can be understood and recovered from.

The pilot is not intended to prove Atlas knows everything.

It is intended to prove that Atlas provides trustworthy and useful infrastructure intelligence from the information it has.

---

# Guiding Statements

> **Truth first. Built for operators.**

> **Build the simplest correct thing.**

> **Parse faithfully first. Correlate later. Derive last.**

> **Infrastructure sources are inputs. Diagrams are outputs.**

> **Unknown is better than invented certainty.**

> **Complexity must earn its place.**

> **Atlas reconstructs the network so the operator does not have to.**

---

# Historical Project

The previous Atlas implementation is retained at:

[Iron-Signal-Systems/old-atlas](https://github.com/Iron-Signal-Systems/old-atlas)

It is historical reference material.

It does not define the architecture or repository structure of the current Atlas implementation.

---

# Security

Atlas is designed to operate with a small, understandable attack surface and strong default trust boundaries.

Core security principles include:

* read-only collection by default;
* least-privilege access to infrastructure sources;
* explicitly scoped collection;
* all source input treated as untrusted;
* bounded parsers and resource usage;
* timeouts where appropriate;
* conservative parsing and validation;
* no secrets or unrestricted customer data committed to Git;
* clear separation between collected facts and derived conclusions;
* explicit handling of unknown, conflicting, or incomplete information;
* straightforward code that can be reviewed and tested;
* dependency and build review;
* tested failure behavior; and
* no silent modification of managed infrastructure.

Atlas should prefer simple, reviewable security boundaries over unnecessary framework complexity.

> **Simple code that can be understood, tested, and reviewed is itself a security property.**

Atlas is currently **pre-alpha and not production ready**.

Security behavior, interfaces, schemas, and implementation details may change before a supported release.

Suspected vulnerabilities should **not** be reported through public GitHub issues.

See [SECURITY.md](SECURITY.md) for reporting instructions and security scope.


---

# License

Atlas is proprietary source-available software.

See [LICENSE](LICENSE) for the applicable terms.

Copyright © 2026 John Joseph Wood. All rights reserved.
