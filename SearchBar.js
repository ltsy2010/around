import React, { useState } from "react";
import { Input, Radio } from "antd";

import { SEARCH_KEY } from "../constants";


const { Search } = Input;

function SearchBar(props) {
    const [searchType, setSearchType] = useState(SEARCH_KEY.all);
    const [error, setError] = useState("");

    const changeSearchType = (e) => {
        const searchType = e.target.value;
        //case1 type = all -> value: ""
        //case2 type == keyword/user -> value: inputvalue => handle search
        if (searchType ===SEARCH_KEY.all) {
            //send search type to home, parent to home
            //fetch all posts
            props.handleSearch({
                type: SEARCH_KEY.all,
                value: ""
            })
        }
        setSearchType(searchType);
        setError(""); //set error with null
    };

    //handle search error
    const handleSearch = (value) => {
        //case 1: display error
        if (searchType !== SEARCH_KEY.all && value === "") {
            setError("Please input your search keyword!");
            return;
        }
        //case2: clear error msg
        setError("");
        //case3:searchtype = keyword/user && value != null
        //send to Home
        props.handleSearch({ type: searchType, keyword: value });
    };

    return (
        <div className="search-bar">
            <Search
                placeholder="input search text"
                enterButton="Search"
                size="large"
                onSearch={handleSearch}
                disabled={searchType === SEARCH_KEY.all}
            />
            <p className="error-msg">{error}</p>

            <Radio.Group
                onChange={changeSearchType}
                value={searchType}
                className="search-type-group"
            >
                <Radio value={SEARCH_KEY.all}>All</Radio>
                <Radio value={SEARCH_KEY.keyword}>Keyword</Radio>
                <Radio value={SEARCH_KEY.user}>User</Radio>
            </Radio.Group>
        </div>
    );
}

export default SearchBar;



