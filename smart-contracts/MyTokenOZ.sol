pragma solidity ^0.6.10;

import "github.com/OpenZeppelin/openzeppelin-contracts/contracts/token/ERC20/ERC20.sol";

contract MyTokenOZ is ERC20 {
    constructor(string memory _name, string memory _symbol) public
    ERC20(_name,_symbol) 
    {
        
    }
    
    function mint(address _recipient, uint _amount) public {
        _mint(_recipient,_amount);
    }
}
