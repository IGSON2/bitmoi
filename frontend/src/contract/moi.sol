// SPDX-License-Identifier: LGPL-3.0-only

pragma solidity ^0.8.1;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract MoiToken is Ownable, ERC20 {
    bool public limited;
    uint256 public maxHoldingAmount;
    uint256 public minHoldingAmount;
    address public uniswapV2Pair = address(1);
    mapping(address => bool) public blacklists;

    mapping(string => bool) public adLocations;
    string[] public locationKeys;

    mapping(string => mapping(address => uint256)) private locks;
    mapping(string => address[]) private lockAddress;

    struct sortedAddrBalance {
        address addr;
        uint256 amount;
    }

    event TokensLocked(
        address indexed account,
        uint256 amount,
        string location
    );
    event TokensUnlocked(
        address indexed account,
        uint256 amount,
        string location
    );
    event LocationAdded(string location);
    event LocationRemoved(string location);

    constructor(uint256 _totalSupply) ERC20("Bitmoi", "MOI") {
        _mint(msg.sender, _totalSupply);
    }

    modifier mustExistLoc(string memory _adLocation) {
        require(adLocations[_adLocation], "Incorrect AD location");
        _;
    }

    modifier mustNotExistLoc(string memory _adLocation) {
        require(!adLocations[_adLocation], "Location is already exist");
        _;
    }

    modifier moreThanZeroAmt(uint256 _amount) {
        require(_amount > 0, "Amount must be greater than zero");
        _;
    }

    function setAdLocation(
        string memory _adLocation
    ) external onlyOwner mustNotExistLoc(_adLocation) {
        adLocations[_adLocation] = true;
        locationKeys.push(_adLocation);
        emit LocationAdded(_adLocation);
    }

    function removeAdLocation(
        string memory _adLocation
    ) external onlyOwner mustExistLoc(_adLocation) {
        delete adLocations[_adLocation];

        uint256 idx = findLocationIdx(_adLocation);
        require(idx > 0, "Invalid location for location key array");
        delete locationKeys[idx];

        emit LocationRemoved(_adLocation);
    }

    function findLocationIdx(
        string memory _adLocation
    ) private view returns (uint256) {
        uint256 length = locationKeys.length;

        for (uint256 i = 0; i < length; i++) {
            if (
                keccak256(bytes(locationKeys[i])) ==
                keccak256(bytes(_adLocation))
            ) {
                return i;
            }
        }

        return 0;
    }

    function lockTokens(
        address _from,
        uint256 _amount,
        string memory _adLocation
    ) external onlyOwner mustExistLoc(_adLocation) moreThanZeroAmt(_amount) {
        _transfer(_from, address(this), _amount);

        locks[_adLocation][_from] += _amount;

        lockAddress[_adLocation].push(_from);

        emit TokensLocked(_from, _amount, _adLocation);
    }

    function unlockTokens() external onlyOwner {
        for (uint256 i = 0; i < locationKeys.length; i++) {
            string memory loc = locationKeys[i];
            for (uint256 j = 0; j < lockAddress[loc].length; j++) {
                address addr = lockAddress[loc][j];
                uint256 transferAmt = locks[loc][addr];
                if (transferAmt == 0) {
                    continue;
                }
                locks[loc][addr] = 0;
                _transfer(address(this), addr, transferAmt);
                emit TokensUnlocked(addr, transferAmt, loc);
                delete locks[loc][addr];
            }
            delete locationKeys[i];
        }
    }

    function getLockedTokens(
        address _from,
        string memory _adLocation
    ) public view returns (uint256) {
        uint256 lockedAmount = locks[_adLocation][_from];

        return lockedAmount;
    }

    function getLocksByLocation(
        string memory _adLocation
    ) public view returns (uint256) {
        uint256 addrCnt;
        uint256 length = lockAddress[_adLocation].length;

        for (uint256 i = 0; i < length; i++) {
            address _addr = lockAddress[_adLocation][i];
            if (locks[_adLocation][_addr] > 0) {
                addrCnt++;
            }
        }
        return addrCnt;
    }

    function sortByLocation(
        string memory _adLocation
    ) public view returns (sortedAddrBalance[] memory) {
        sortedAddrBalance[] memory sortedAddr;

        for (uint256 i = 0; i < lockAddress[_adLocation].length; i++) {
            address _addr = lockAddress[_adLocation][i];
            uint256 _amount = locks[_adLocation][_addr];
            for (uint256 j = 0; j < sortedAddr.length; j++) {
                if (_amount > sortedAddr[j].amount) {
                    (sortedAddr[j], sortedAddr[j + 1]) = (
                        sortedAddrBalance(_addr, _amount),
                        sortedAddr[j]
                    );
                } else if (_amount == sortedAddr[j].amount) {
                    (sortedAddr[j + 1], sortedAddr[j]) = (
                        sortedAddr[j],
                        sortedAddrBalance(_addr, _amount)
                    ); // 먼저 예치한 사람이 우위
                }
                sortedAddr[0] = (sortedAddrBalance(_addr, _amount));
            }
        }

        return sortedAddr;
    }

    function blacklist(
        address _address,
        bool _isBlacklisting
    ) external onlyOwner {
        blacklists[_address] = _isBlacklisting;
    }

    function setRule(
        bool _limited,
        address _uniswapV2Pair,
        uint256 _maxHoldingAmount,
        uint256 _minHoldingAmount
    ) external onlyOwner {
        limited = _limited;
        uniswapV2Pair = _uniswapV2Pair;
        maxHoldingAmount = _maxHoldingAmount;
        minHoldingAmount = _minHoldingAmount;
    }

    function _beforeTokenTransfer(
        address from,
        address to,
        uint256 amount
    ) internal virtual override {
        require(!blacklists[to] && !blacklists[from], "Blacklisted");

        if (uniswapV2Pair == address(0)) {
            require(from == owner() || to == owner(), "trading is not started");
            return;
        }

        if (limited && from == uniswapV2Pair) {
            require(
                super.balanceOf(to) + amount <= maxHoldingAmount &&
                    super.balanceOf(to) + amount >= minHoldingAmount,
                "Forbid"
            );
        }
    }

    function burn(address _addr, uint256 value) external onlyOwner {
        _burn(_addr, value);
    }

    function decimals() public view virtual override returns (uint8) {
        return 1;
    }
}
