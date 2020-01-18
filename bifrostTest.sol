/**
 *Submitted for verification at Etherscan.io on 2020-01-15
*/

pragma solidity >=0.4.24;
contract BifrostTest {
	mapping(address => uint256) store;
	function setValue(uint256 value) public {
		store[msg.sender] = value;
	}
	function getValue() public view returns(uint256) {
		return store[msg.sender];
	}
}