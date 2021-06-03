import React, { useState, useEffect } from "react";
import { Tabs, message, Row, Col, Button } from "antd";
import axios from "axios";
//use tabs from antd
import SearchBar from "./SearchBar";
import PhotoGallery from "./PhotoGallery";
import { SEARCH_KEY, BASE_URL, TOKEN_KEY } from "../constants";
import CreatePostButton from "./CreatePostButton";

const { TabPane } = Tabs;

function Home(props) {
    const [posts, setPost] = useState([]); //initial is empty
    const [activeTab, setActiveTab] = useState("image");
    const [searchOption, setSearchOption] = useState({
        type: SEARCH_KEY.all,
        keyword: ""
    });

    //when search -> call useeffect
    //when searchoption changes, fetch data from server -> do useEffect
    //deps: array -> multiple parameters.
    //do search first time -> didmount -> search: {type: search_key.all, keyword: ""}
    //do search after the first time -> didupdate -> search: {type: keyword/user, keyword: value}
    useEffect(() => {
        const { type, keyword } = searchOption;
        fetchPost(searchOption);
    }, [searchOption]);

    //option: parameter, option url
    const fetchPost = (option) => {
        const { type, keyword } = option;
        let url = "";

        if (type === SEARCH_KEY.all) {
            url = `${BASE_URL}/search`;
        } else if (type === SEARCH_KEY.user) {
            url = `${BASE_URL}/search?user=${keyword}`;
        } else {
            url = `${BASE_URL}/search?keywords=${keyword}`;
        }

        //fetch data from server
        const opt = {
            method: "GET",
            url: url,
            headers: {
                Authorization: `Bearer ${localStorage.getItem(TOKEN_KEY)}`
            }
        };

        axios(opt)
            .then((res) => {
                if (res.status === 200) {
                    //setstate -> post ->rerender
                    setPost(res.data);
                }
            })
            .catch((err) => {
                message.error("Fetch posts failed!");
                console.log("fetch posts failed: ", err.message);
            });
    };

    const renderPosts = (type) => {
        //case1: no posts => display no data
        if (!posts || posts.length === 0) {
            return <div>No data!</div>;
        }
        //type = image, display images
        if (type === "image") {
            //filter image from post
            //map: 遍历image item
            const imageArr = posts
                .filter((item) => item.type === "image")
                .map((image) => {
                    return {
                        src: image.url,
                        user: image.user,
                        caption: image.message,
                        thumbnail: image.url,
                        thumbnailWidth: 300,
                        thumbnailHeight: 200,
                        postId: image.id
                    };
                });

            console.log("images -> ", posts);
            return <PhotoGallery images={imageArr} />;
        } else if (type === "video") {
            return (
                <Row gutter={32}>
                    {posts
                        .filter((post) => post.type === "video")
                        .map((post) => (
                            <Col span={8} key={post.url}>
                                <video src={post.url} controls={true} className="video-block" />
                                <p>
                                    {post.user}: {post.message}
                                </p>
                            </Col>
                        ))}
                </Row>
            );

        }
    };


    const handleSearch = (option) => {
        const { type, keyword } = option;
        setSearchOption({ type: type, keyword: keyword });
    };

    //显示拿到的数据 type: image/video
    //用settimeout when there is no notification
    const showPost = (type) => {
        console.log("type -> ", type);
        setActiveTab(type);

        setTimeout(() => {
            setSearchOption({ type: SEARCH_KEY.all, keyword: "" });
        }, 3000);
    };

    const operations = <CreatePostButton onShowPost={showPost} />;

    return (
        <div className="home">
            <SearchBar handleSearch={handleSearch} />
            <div className="display">
                <Tabs
                    onChange={(key) => setActiveTab(key)}
                    defaultActiveKey="image"  //default is image
                    activeKey={activeTab} //change, not fixed
                    tabBarExtraContent={operations}
                >
                    <TabPane tab="Images" key="image">
                        {renderPosts("image")}
                    </TabPane>
                    <TabPane tab="Videos" key="video">
                        {renderPosts("video")}
                    </TabPane>
                </Tabs>
            </div>
        </div>
    );
}

export default Home;


