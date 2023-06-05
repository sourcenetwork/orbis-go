// Package rabin implements the protocol described in
// "Secure Distributed Key Generation for Discrete-Log
// Based Cryptosystems" by R. Gennaro, S. Jarecki, H. Krawczyk, and T. Rabin.
// DKG enables a group of participants to generate a distributed key
// with each participants holding only a share of the key. The key is also
// never computed locally but generated distributively whereas the public part
// of the key is known by every participants.
// The underlying basis for this protocol is the VSS protocol implemented in the
// share/vss package.
//
// The protocol works as follow:
//
//   1. Each participant instantiates a DistKeyShare (DKS) struct.
//   2. Then each participant runs an instance of the VSS protocol:
//     - each participant generates their deals with the method `Deals()` and then
//      sends them to the right recipient.
//     - each participant processes the received deal with `ProcessDeal()` and
//      broadcasts the resulting response.
//     - each participant processes the response with `ProcessResponse()`. If a
//      justification is returned, it must be broadcasted.
//   3. Each participant can check if step 2. is done by calling
//   `Certified()`.Those participants where Certified() returned true, belong to
//   the set of "qualified" participants who will generate the distributed
//   secret. To get the list of qualified participants, use QUAL().
//   4. Each QUAL participant generates their secret commitments calling
//    `SecretCommits()` and broadcasts them to the QUAL set.
//   5. Each QUAL participant processes the received secret commitments using
//    `ProcessSecretCommits()`. If there is an error, it can return a commitment complaint
//    (ComplaintCommits) that must be broadcasted to the QUAL set.
//   6. Each QUAL participant receiving a complaint can process it with
//    `ProcessComplaintCommits()` which returns the secret share
//    (ReconstructCommits) given from the malicious participant. This structure
//    must be broadcasted to all the QUAL participant.
//   7. At this point, every QUAL participant can issue the distributed key by
//    calling `DistKeyShare()`.

package rabin
