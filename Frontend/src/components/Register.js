import React from "react";
import { Form, Input, Button, message, Select, InputNumber } from 'antd';
import axios from 'axios';

import { BASE_URL } from "../constants";

const formItemLayout = {
    labelCol: {
        xs: { span: 24 },
        sm: { span: 8 },
    },
    wrapperCol: {
        xs: { span: 24 },
        sm: { span: 16 },
    },
};
const tailFormItemLayout = {
    wrapperCol: {
        xs: {
            span: 16,
            offset: 0,
        },
        sm: {
            span: 16,
            offset: 8,
        },
    },
};
const { Option } = Select;

function Register(props) {
    const [form] = Form.useForm();

    const onFinish = values => {
        console.log('Received values of form: ', values);
        const { username, password, age, gender } = values;
        const opt = {
            method: 'POST',
            url: `${BASE_URL}/signup`,
            data: {
                username: username,
                password: password,
                age: age,
                gender: gender,
            },
            headers: { 'content-type': 'application/json' }
        };

        axios(opt)
            .then(response => {
                console.log(response)
                // case1: registered success
                if (response.status === 200) {
                    message.success('Registration succeed!');
                    props.history.push('/login');
                }
            })
            .catch(error => {
                console.log('register failed: ', error.message);
                message.success('Registration failed!');
                // throw new Error('Signup Failed!')
            })
    };

    return (
        <Form
            {...formItemLayout}
            form={form}
            name="register"
            onFinish={onFinish}
            className="register"
        >
            <Form.Item
                name="username"
                label="Username"
                rules={[
                    {
                        required: true,
                        message: 'Please input your Username!',
                    },
                ]}
            >
                <Input />
            </Form.Item>

            <Form.Item
                name="age"
                label="Age"
                rules={[
                    {
                        required: true,
                        type: 'number',
                        min: 0,
                        max: 120,
                    },
                ]}
            >
                <InputNumber className="register-age" />
            </Form.Item>

            <Form.Item 
                name="gender" 
                label="Gender" 
                rules={[
                    { 
                        required: true,
                        message: "Please input your age!",
                    }
                ]}
            >
                <Select
                    placeholder="Select a option and change input text above"
                    allowClear
                >
                    <Option value="male">male</Option>
                    <Option value="female">female</Option>
                    <Option value="other">other</Option>
                </Select>
            </Form.Item>

            <Form.Item
                name="password"
                label="Password"
                rules={[
                    {
                        required: true,
                        message: 'Please input your password!',
                    },
                ]}
                hasFeedback
            >
                <Input.Password />
            </Form.Item>

            <Form.Item
                name="confirm"
                label="Confirm Password"
                dependencies={['password']}
                hasFeedback
                rules={[
                    {
                        required: true,
                        message: 'Please confirm your password!',
                    },
                    ({ getFieldValue }) => ({
                        validator(rule, value) {
                            if (!value || getFieldValue('password') === value) {
                                return Promise.resolve();
                            }
                            return Promise.reject('The two passwords that you entered do not match!');
                        },
                    }),
                ]}
            >
                <Input.Password />
            </Form.Item>

            <Form.Item {...tailFormItemLayout}>
                <Button type="primary" htmlType="submit" className="register-btn">
                    Register
                </Button>
            </Form.Item>
        </Form>
    );
}

export default Register;