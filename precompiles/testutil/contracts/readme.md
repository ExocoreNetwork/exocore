# how to write a test to call the precompile contract from a contract

## cmd to generate abi and bin

`solc --base-path ./ --include-path ./../.. --evm-version paris --bin --abi ./DepositCaller.sol -o . --overwrite`

First you need to create a file named DepositCaller.json and add the generated bin and abi to the created json file.Then
you can write some tests to call the Deposit precompile contract from contract account. You can refer to the file
`deposit_integrate_test.go` to get how to write the test codes.
