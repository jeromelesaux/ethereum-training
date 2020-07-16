pragma solidity ^0.6.10;


contract MyToken {
    
    // variables
    // balances prend en cle et retourne la balance de l'adresse en valeur
    // balances[0x...] = 10 si 0x.... détient 10 tokens
    mapping(address => uint256) balances;
    
    string private  name;
    string public symbol;
    uint8 public decimals;
    
    uint256 private _totalSupply;
    
    // proprietaire qui va pouvoir créer des tokens
    address owner;
    
    // evenements 
    event Transfer(address sender, address receiver, uint256 amount);
    
    // constructeur
    constructor(string memory _name, string memory _symbol, uint8 _decimals, uint256 _newTotalSupply) public {
        name = _name;
        symbol = _symbol;
        decimals = _decimals;
        _totalSupply = _newTotalSupply;
        //designe le sender comme proprietaire
        owner = msg.sender;
        
        // affecte le supply au proprietaire
        balances[msg.sender] = _newTotalSupply;
        
        emit Transfer(address(0), msg.sender,_newTotalSupply);
    }
    
    
    // methodes
    // tranfert 
    function transfer(address _receiver, uint256 _amount) public {
        // requires : balance de l'utilisateur est supérieur ou égale au montant transféré
        require(balances[msg.sender] >= _amount);
        
        // logique 
        balances[msg.sender] -= _amount;
        balances[_receiver] += _amount;
        
        // evenement 
        emit Transfer(msg.sender,_receiver,_amount);
    }
    
    
    // creation de tokens
    function mint(address _receiver, uint256 _amount) public {
        // require : est ce que le sender a le droit de créer des tokens 
        require(msg.sender == owner, "sender is not the contract owner");
        
        // logique 
        balances[_receiver] += _amount;
        _totalSupply += _amount;
        
        // evenement 
        emit Transfer(address(0),_receiver,_amount);
        
    }
    
    // Bruler des tokens
    function burn(uint256 _amount) public {
        // require : verifier que l'utilisateur a suffisament de tokens
        require(balances[msg.sender] >= _amount);
        
        balances[msg.sender] -= _amount;
        _totalSupply -= _amount;
        
        emit Transfer(msg.sender,address(0),_amount);
        
    }
    
    // getter de la balance d'un utilisateur 
    function balanceOf(address _user) public view returns( uint256) {
       return balances[_user];
    }
    
    function totalSupply() public view returns(uint256){
        return _totalSupply;
    }
}

