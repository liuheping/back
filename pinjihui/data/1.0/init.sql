--
-- PostgreSQL database dump
--

-- Dumped from database version 10.3 (Ubuntu 10.3-1)
-- Dumped by pg_dump version 10.3 (Ubuntu 10.3-1)

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: address; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.address AS (
	area_id integer,
	region_name character varying(255),
	address character varying(255)
);


ALTER TYPE public.address OWNER TO postgres;

--
-- Name: brand_type; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.brand_type AS ENUM (
    'part',
    'excavator'
);


ALTER TYPE public.brand_type OWNER TO postgres;

--
-- Name: cash_request_status; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.cash_request_status AS ENUM (
    'unchecked',
    'checking',
    'checked',
    'paid',
    'finished',
    'refused',
    'closed'
);


ALTER TYPE public.cash_request_status OWNER TO postgres;

--
-- Name: debit_card; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.debit_card AS (
	card_holder character varying(255),
	bank_name character varying(255),
	card_number character varying(32),
	province integer,
	city integer,
	branch character varying(255)
);


ALTER TYPE public.debit_card OWNER TO postgres;

--
-- Name: favorite_type; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.favorite_type AS ENUM (
    'product',
    'store'
);


ALTER TYPE public.favorite_type OWNER TO postgres;

--
-- Name: how_oos; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.how_oos AS ENUM (
    'together',
    'cancel',
    'consult'
);


ALTER TYPE public.how_oos OWNER TO postgres;

--
-- Name: input_types; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.input_types AS ENUM (
    'textfield',
    'textarea',
    'dropdown',
    'mutiselect',
    'time'
);


ALTER TYPE public.input_types OWNER TO postgres;

--
-- Name: order_status; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.order_status AS ENUM (
    'unconfirmed',
    'confirmed',
    'cancelled',
    'invalid',
    'returned'
);


ALTER TYPE public.order_status OWNER TO postgres;

--
-- Name: pay_status; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.pay_status AS ENUM (
    'unpaid',
    'paying',
    'paid'
);


ALTER TYPE public.pay_status OWNER TO postgres;

--
-- Name: shipping_address; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.shipping_address AS (
	consignee character varying(255),
	address public.address,
	zipcode character varying(32),
	mobile character varying(32)
);


ALTER TYPE public.shipping_address OWNER TO postgres;

--
-- Name: shipping_status; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.shipping_status AS ENUM (
    'unshipped',
    'shipped',
    'invalid',
    'returned'
);


ALTER TYPE public.shipping_status OWNER TO postgres;

--
-- Name: user_status; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.user_status AS ENUM (
    'normal',
    'banned',
    'unchecked',
    'checked'
);


ALTER TYPE public.user_status OWNER TO postgres;

--
-- Name: user_type; Type: TYPE; Schema: public; Owner: postgres
--

CREATE TYPE public.user_type AS ENUM (
    'admin',
    'consumer',
    'provider',
    'ally'
);


ALTER TYPE public.user_type OWNER TO postgres;

--
-- Name: auto_update_timestamp(); Type: FUNCTION; Schema: public; Owner: postgres
--

CREATE FUNCTION public.auto_update_timestamp() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
begin
  new.updated_at = current_timestamp;
  return new;
end;
$$;


ALTER FUNCTION public.auto_update_timestamp() OWNER TO postgres;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: addresses; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.addresses (
    id character varying(32) NOT NULL,
    user_id character varying(32) NOT NULL,
    consignee character varying(32) NOT NULL,
    address public.address NOT NULL,
    zipcode character varying(6),
    mobile character varying(11),
    is_default boolean DEFAULT false NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public.addresses OWNER TO postgres;

--
-- Name: TABLE addresses; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.addresses IS '用户收货地址表';


--
-- Name: COLUMN addresses.consignee; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.addresses.consignee IS '收货人姓名';


--
-- Name: attribute_sets; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.attribute_sets (
    id character varying(32) NOT NULL,
    name character varying(255) NOT NULL,
    attribute_ids character varying[]
);


ALTER TABLE public.attribute_sets OWNER TO postgres;

--
-- Name: TABLE attribute_sets; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.attribute_sets IS '属性集表';


--
-- Name: COLUMN attribute_sets.attribute_ids; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.attribute_sets.attribute_ids IS '属性ID';


--
-- Name: attributes; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.attributes (
    id character varying(255) NOT NULL,
    name character varying(255) NOT NULL,
    type public.input_types NOT NULL,
    required boolean DEFAULT false NOT NULL,
    default_value character varying(255),
    options character varying[],
    merchant_id character varying(32),
    enabled boolean DEFAULT true NOT NULL,
    searchable boolean DEFAULT false NOT NULL,
    deleted boolean DEFAULT false NOT NULL,
    code character varying(32) NOT NULL
);


ALTER TABLE public.attributes OWNER TO postgres;

--
-- Name: TABLE attributes; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.attributes IS '属性表';


--
-- Name: COLUMN attributes.name; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.attributes.name IS '名字';


--
-- Name: COLUMN attributes.type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.attributes.type IS '输入类型';


--
-- Name: COLUMN attributes.required; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.attributes.required IS '是否必须';


--
-- Name: COLUMN attributes.default_value; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.attributes.default_value IS '默认值';


--
-- Name: COLUMN attributes.options; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.attributes.options IS '选项值';


--
-- Name: COLUMN attributes.merchant_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.attributes.merchant_id IS '商家ID，公共的为空';


--
-- Name: COLUMN attributes.enabled; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.attributes.enabled IS '是否启用';


--
-- Name: COLUMN attributes.searchable; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.attributes.searchable IS '是否可搜索';


--
-- Name: COLUMN attributes.deleted; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.attributes.deleted IS '是否删除';


--
-- Name: COLUMN attributes.code; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.attributes.code IS '属性代码';


--
-- Name: brands; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.brands (
    id character varying(32) NOT NULL,
    name character varying(255),
    thumbnail character varying(255),
    description text,
    deleted boolean DEFAULT false NOT NULL,
    enabled boolean DEFAULT true NOT NULL,
    sort_order integer DEFAULT 255 NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    type public.brand_type NOT NULL,
    machine_types character varying[]
);


ALTER TABLE public.brands OWNER TO postgres;

--
-- Name: TABLE brands; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.brands IS '品牌表';


--
-- Name: COLUMN brands.thumbnail; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.brands.thumbnail IS '缩略图';


--
-- Name: COLUMN brands.type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.brands.type IS '品牌类型,有挖机品牌和部件品牌';


--
-- Name: COLUMN brands.machine_types; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.brands.machine_types IS '此品牌包含的机型';


--
-- Name: carts; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.carts (
    id character varying(32) NOT NULL,
    product_id character varying(32) NOT NULL,
    user_id character varying(32) NOT NULL,
    product_count integer NOT NULL,
    merchant_id character varying(32) NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public.carts OWNER TO postgres;

--
-- Name: TABLE carts; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.carts IS '购物车';


--
-- Name: COLUMN carts.product_count; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.carts.product_count IS '商品数量';


--
-- Name: COLUMN carts.merchant_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.carts.merchant_id IS '商家ID';


--
-- Name: cash_request; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.cash_request (
    id character varying(32) NOT NULL,
    amount money NOT NULL,
    debit_card_info public.debit_card NOT NULL,
    status public.cash_request_status DEFAULT 'unchecked'::public.cash_request_status NOT NULL,
    reply text,
    note text
);


ALTER TABLE public.cash_request OWNER TO postgres;

--
-- Name: categories; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.categories (
    id character varying(32) NOT NULL,
    parent_id character varying(32),
    name character varying(255) NOT NULL,
    sort_order integer DEFAULT 255,
    thumbnail character varying(255),
    deleted boolean DEFAULT false NOT NULL,
    enabled boolean DEFAULT true NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public.categories OWNER TO postgres;

--
-- Name: TABLE categories; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.categories IS '分类表';


--
-- Name: COLUMN categories.sort_order; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.categories.sort_order IS '排序号';


--
-- Name: COLUMN categories.thumbnail; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.categories.thumbnail IS '缩略图';


--
-- Name: comments; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.comments (
    id character varying(32) NOT NULL,
    user_id character varying(32) NOT NULL,
    product_id character varying(32) NOT NULL,
    rank smallint,
    order_id character varying(32) NOT NULL,
    content text NOT NULL,
    is_show boolean DEFAULT false NOT NULL,
    created_at timestamp(6) without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    user_ip inet,
    reply text,
    merchant_id character varying(32) NOT NULL
);


ALTER TABLE public.comments OWNER TO postgres;

--
-- Name: TABLE comments; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.comments IS '评价表';


--
-- Name: COLUMN comments.rank; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.comments.rank IS '星级;只有1 到5 星;由数字代替;其中5 代表5 星';


--
-- Name: COLUMN comments.is_show; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.comments.is_show IS '是否审核通过（通过才显示）';


--
-- Name: COLUMN comments.user_ip; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.comments.user_ip IS '用户评论时IP';


--
-- Name: COLUMN comments.reply; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.comments.reply IS '卖家回复';


--
-- Name: configs; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.configs (
    id character varying(32) NOT NULL,
    code character varying(32) NOT NULL,
    value text NOT NULL,
    sort_order smallint DEFAULT 255,
    name character varying(255) NOT NULL,
    description text,
    deleted boolean DEFAULT false NOT NULL
);


ALTER TABLE public.configs OWNER TO postgres;

--
-- Name: TABLE configs; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.configs IS '全站配置表';


--
-- Name: COLUMN configs.code; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.configs.code IS '跟变量名的作用差不多，其实就是语言包中的字符串索引，如$_LANG[''''cfg_range''''][''''cart_confirm'''']';


--
-- Name: debit_card_info; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.debit_card_info (
    user_id character varying(32) NOT NULL,
    debit_card_info public.debit_card NOT NULL,
    is_checked boolean DEFAULT false NOT NULL,
    created_at timestamp(6) without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone
);


ALTER TABLE public.debit_card_info OWNER TO postgres;

--
-- Name: COLUMN debit_card_info.user_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.debit_card_info.user_id IS '用户ID';


--
-- Name: COLUMN debit_card_info.debit_card_info; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.debit_card_info.debit_card_info IS '银行卡信息';


--
-- Name: COLUMN debit_card_info.is_checked; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.debit_card_info.is_checked IS '是否审核';


--
-- Name: COLUMN debit_card_info.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.debit_card_info.created_at IS '创建时间';


--
-- Name: COLUMN debit_card_info.updated_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.debit_card_info.updated_at IS '更新时间';


--
-- Name: favorites; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.favorites (
    id character varying(32) NOT NULL,
    user_id character varying(32) NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    type public.favorite_type NOT NULL,
    object_id character varying(32) NOT NULL
);


ALTER TABLE public.favorites OWNER TO postgres;

--
-- Name: TABLE favorites; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.favorites IS '收藏表';


--
-- Name: COLUMN favorites.type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.favorites.type IS '收藏对象的类型';


--
-- Name: goose_db_version; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.goose_db_version (
    id integer NOT NULL,
    version_id bigint NOT NULL,
    is_applied boolean NOT NULL,
    tstamp timestamp without time zone DEFAULT now()
);


ALTER TABLE public.goose_db_version OWNER TO postgres;

--
-- Name: goose_db_version_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.goose_db_version_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.goose_db_version_id_seq OWNER TO postgres;

--
-- Name: goose_db_version_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.goose_db_version_id_seq OWNED BY public.goose_db_version.id;


--
-- Name: merchant_profiles; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.merchant_profiles (
    user_id character varying(32) NOT NULL,
    social_id character varying(18),
    rep_name character varying(32),
    company_name character varying(255),
    company_address public.address,
    delivery_address public.address,
    license_image character varying(255),
    company_image character varying(255),
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public.merchant_profiles OWNER TO postgres;

--
-- Name: TABLE merchant_profiles; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.merchant_profiles IS '商家信息表';


--
-- Name: COLUMN merchant_profiles.social_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.merchant_profiles.social_id IS '身份证号';


--
-- Name: COLUMN merchant_profiles.rep_name; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.merchant_profiles.rep_name IS '法人姓名/真实姓名';


--
-- Name: COLUMN merchant_profiles.company_name; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.merchant_profiles.company_name IS '公司名称';


--
-- Name: COLUMN merchant_profiles.company_address; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.merchant_profiles.company_address IS '公司地址';


--
-- Name: COLUMN merchant_profiles.delivery_address; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.merchant_profiles.delivery_address IS '发货地址';


--
-- Name: COLUMN merchant_profiles.license_image; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.merchant_profiles.license_image IS '营业执照图片';


--
-- Name: COLUMN merchant_profiles.company_image; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.merchant_profiles.company_image IS '形象照';


--
-- Name: operation_logs; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.operation_logs (
    id character varying(32) NOT NULL,
    user_id character varying(32) NOT NULL,
    action character varying(255) NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL
);


ALTER TABLE public.operation_logs OWNER TO postgres;

--
-- Name: TABLE operation_logs; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.operation_logs IS '操作日志表';


--
-- Name: order_products; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.order_products (
    id character varying(32) NOT NULL,
    order_id character varying(32) NOT NULL,
    product_id character varying(32) NOT NULL,
    product_name character varying(255) NOT NULL,
    product_number smallint NOT NULL,
    product_price money NOT NULL,
    product_image character varying(255) NOT NULL
);


ALTER TABLE public.order_products OWNER TO postgres;

--
-- Name: TABLE order_products; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.order_products IS '商品订单关联表';


--
-- Name: COLUMN order_products.product_name; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.order_products.product_name IS '下单时的商品名';


--
-- Name: COLUMN order_products.product_number; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.order_products.product_number IS '购买数量';


--
-- Name: COLUMN order_products.product_price; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.order_products.product_price IS '下单时商品的售价';


--
-- Name: COLUMN order_products.product_image; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.order_products.product_image IS '商品缩略图';


--
-- Name: orders; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.orders (
    id character varying(32) NOT NULL,
    user_id character varying(32) NOT NULL,
    order_status public.order_status NOT NULL,
    shipping_status public.shipping_status NOT NULL,
    pay_status public.pay_status NOT NULL,
    postscript character varying(255),
    shipping_id character varying(32) NOT NULL,
    shipping_name character varying(255) NOT NULL,
    pay_id character varying(32) NOT NULL,
    pay_name character varying(255) NOT NULL,
    how_oos public.how_oos,
    inv_payee character varying(32),
    inv_type character varying(255),
    inv_content character varying(255),
    amount money NOT NULL,
    shipping_fee money,
    pay_fee money,
    money_paid money,
    order_amount money,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    confirm_time timestamp without time zone,
    pay_time timestamp without time zone,
    shipping_time timestamp without time zone,
    tax money,
    parent_id character varying(32),
    merchant_id character varying(32),
    address public.shipping_address NOT NULL,
    note text
);


ALTER TABLE public.orders OWNER TO postgres;

--
-- Name: TABLE orders; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.orders IS '订单表';


--
-- Name: COLUMN orders.id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.orders.id IS '订单号';


--
-- Name: COLUMN orders.order_status; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.orders.order_status IS '订单状态';


--
-- Name: COLUMN orders.shipping_status; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.orders.shipping_status IS '商品配送情况';


--
-- Name: COLUMN orders.pay_status; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.orders.pay_status IS '支付状态';


--
-- Name: COLUMN orders.postscript; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.orders.postscript IS '订单留言';


--
-- Name: COLUMN orders.shipping_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.orders.shipping_id IS '配送方式ID';


--
-- Name: COLUMN orders.how_oos; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.orders.how_oos IS '缺货/备货处理方式';


--
-- Name: COLUMN orders.inv_payee; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.orders.inv_payee IS '发票抬头';


--
-- Name: COLUMN orders.inv_type; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.orders.inv_type IS '发票类型';


--
-- Name: COLUMN orders.inv_content; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.orders.inv_content IS '发票内容，取值配置表';


--
-- Name: COLUMN orders.amount; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.orders.amount IS '商品总价格';


--
-- Name: COLUMN orders.shipping_fee; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.orders.shipping_fee IS '配送费用';


--
-- Name: COLUMN orders.pay_fee; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.orders.pay_fee IS '支付费用,跟支付方式的配置相关,取值表payment';


--
-- Name: COLUMN orders.money_paid; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.orders.money_paid IS '已付款金额';


--
-- Name: COLUMN orders.order_amount; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.orders.order_amount IS '应付款金额';


--
-- Name: COLUMN orders.created_at; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.orders.created_at IS '订单生成时间';


--
-- Name: COLUMN orders.confirm_time; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.orders.confirm_time IS '订单确认时间';


--
-- Name: COLUMN orders.pay_time; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.orders.pay_time IS '支付时间';


--
-- Name: COLUMN orders.shipping_time; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.orders.shipping_time IS '订单配送时间';


--
-- Name: COLUMN orders.tax; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.orders.tax IS '税额';


--
-- Name: COLUMN orders.merchant_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.orders.merchant_id IS '商家ID';


--
-- Name: payments; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.payments (
    id character varying(32) NOT NULL,
    pay_name character varying(255) NOT NULL,
    pay_code character varying(32) NOT NULL,
    pay_fee text,
    pay_desc text,
    sort_order integer DEFAULT 255 NOT NULL,
    enabled boolean DEFAULT true NOT NULL,
    is_cod boolean DEFAULT false NOT NULL,
    is_online boolean DEFAULT true NOT NULL,
    deleted boolean DEFAULT false NOT NULL
);


ALTER TABLE public.payments OWNER TO postgres;

--
-- Name: TABLE payments; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.payments IS '支付方式表';


--
-- Name: COLUMN payments.pay_code; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.payments.pay_code IS '支付方式的英文缩写,其实是该支付方式处理插件的不带后缀的文件名部分';


--
-- Name: COLUMN payments.pay_fee; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.payments.pay_fee IS '支付费用';


--
-- Name: COLUMN payments.pay_desc; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.payments.pay_desc IS '支付方式描述';


--
-- Name: COLUMN payments.sort_order; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.payments.sort_order IS '支付方式的显示顺序';


--
-- Name: COLUMN payments.enabled; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.payments.enabled IS '是否可用';


--
-- Name: COLUMN payments.is_cod; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.payments.is_cod IS '是否货到付款';


--
-- Name: COLUMN payments.is_online; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.payments.is_online IS '是否在线支付';


--
-- Name: COLUMN payments.deleted; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.payments.deleted IS '是否删除';


--
-- Name: product_images; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.product_images (
    id character varying(32) NOT NULL,
    product_id character varying(32) NOT NULL,
    small_image character varying(255),
    medium_image character varying(255),
    big_image character varying(255) NOT NULL
);


ALTER TABLE public.product_images OWNER TO postgres;

--
-- Name: TABLE product_images; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.product_images IS '商品图片表';


--
-- Name: COLUMN product_images.small_image; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.product_images.small_image IS '缩略图';


--
-- Name: COLUMN product_images.medium_image; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.product_images.medium_image IS '中等大小图片';


--
-- Name: COLUMN product_images.big_image; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.product_images.big_image IS '大图（原图）';


--
-- Name: products; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.products (
    id character varying(32) NOT NULL,
    name character varying(255) NOT NULL,
    is_sale boolean DEFAULT false NOT NULL,
    attribute_set_id character varying(32),
    batch_price numeric(100,2) NOT NULL,
    second_price numeric(100,2) NOT NULL,
    category_id character varying(32),
    related_ids character varying[],
    content text DEFAULT ''::text,
    brand_id character varying(32),
    deleted boolean DEFAULT false NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    tags character varying[],
    attrs jsonb,
    recommended boolean DEFAULT false NOT NULL,
    sales_volume integer DEFAULT 0 NOT NULL,
    machine_types character varying[],
    retail_price     numeric
);


ALTER TABLE public.products OWNER TO postgres;

--
-- Name: TABLE products; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.products IS '商品表';


--
-- Name: COLUMN products.batch_price; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.products.batch_price IS '批发价,即供货商给平台的价格';


--
-- Name: COLUMN products.second_price; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.products.second_price IS '平台给加盟商的价格';


--
-- Name: COLUMN products.related_ids; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.products.related_ids IS '关联商品ID';


--
-- Name: COLUMN products.content; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.products.content IS '商品详情';


--
-- Name: COLUMN products.brand_id; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.products.brand_id IS '品牌';


--
-- Name: COLUMN products.tags; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.products.tags IS '标签';


--
-- Name: COLUMN products.attrs; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.products.attrs IS '属性';


--
-- Name: COLUMN products.recommended; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.products.recommended IS '是否推荐';


--
-- Name: COLUMN products.machine_types; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.products.machine_types IS '产品适配机型';


--
-- Name: regions; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.regions (
    id integer NOT NULL,
    parent_id integer DEFAULT 0 NOT NULL,
    name character varying(32) NOT NULL,
    sort_order integer DEFAULT 255 NOT NULL
);


ALTER TABLE public.regions OWNER TO postgres;

--
-- Name: TABLE regions; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.regions IS '行政区域表';


--
-- Name: regions_id_seq; Type: SEQUENCE; Schema: public; Owner: postgres
--

CREATE SEQUENCE public.regions_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE public.regions_id_seq OWNER TO postgres;

--
-- Name: regions_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: postgres
--

ALTER SEQUENCE public.regions_id_seq OWNED BY public.regions.id;


--
-- Name: rel_merchants_products; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.rel_merchants_products (
    product_id character varying(32) NOT NULL,
    merchant_id character varying(32) NOT NULL,
    stock integer DEFAULT 1 NOT NULL,
    retail_price numeric(100,2)
);


ALTER TABLE public.rel_merchants_products OWNER TO postgres;

--
-- Name: TABLE rel_merchants_products; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.rel_merchants_products IS '商家商品关联表';


--
-- Name: COLUMN rel_merchants_products.stock; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.rel_merchants_products.stock IS '库存';


--
-- Name: COLUMN rel_merchants_products.retail_price; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.rel_merchants_products.retail_price IS '最终售价';


--
-- Name: users; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.users (
    id character varying(32) NOT NULL,
    name character varying(32),
    mobile character varying(32),
    password bytea NOT NULL,
    type public.user_type NOT NULL,
    email character varying(255),
    created_at timestamp(6) without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp(6) without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    status public.user_status DEFAULT 'normal'::public.user_status NOT NULL,
    last_ip inet,
    last_login_time timestamp(6) without time zone
);


ALTER TABLE public.users OWNER TO postgres;

--
-- Name: TABLE users; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.users IS '用户表';


--
-- Name: COLUMN users.status; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.users.status IS '用户状态,normal:正常;banned:封禁;unchecked: 未审核;checked: 审核通过';


--
-- Name: COLUMN users.last_login_time; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.users.last_login_time IS '最近一次登录时间';


--
-- Name: wecharts; Type: TABLE; Schema: public; Owner: postgres
--

CREATE TABLE public.wecharts (
    openid character varying(32) NOT NULL,
    session_key character varying(32) NOT NULL,
    nick_name character varying(255) NOT NULL,
    gender smallint DEFAULT 0 NOT NULL,
    language character varying(32),
    city character varying(32),
    province character varying(32),
    country character varying(32),
    avatar_url character varying(526),
    user_id character varying(32)
);


ALTER TABLE public.wecharts OWNER TO postgres;

--
-- Name: TABLE wecharts; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON TABLE public.wecharts IS '微信账号信息表';


--
-- Name: COLUMN wecharts.gender; Type: COMMENT; Schema: public; Owner: postgres
--

COMMENT ON COLUMN public.wecharts.gender IS '1:男,2:女,0未知';


--
-- Name: goose_db_version id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.goose_db_version ALTER COLUMN id SET DEFAULT nextval('public.goose_db_version_id_seq'::regclass);


--
-- Name: regions id; Type: DEFAULT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.regions ALTER COLUMN id SET DEFAULT nextval('public.regions_id_seq'::regclass);


--
-- Data for Name: addresses; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.addresses (id, user_id, consignee, address, zipcode, mobile, is_default, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: attribute_sets; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.attribute_sets (id, name, attribute_ids) FROM stdin;
\.


--
-- Data for Name: attributes; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.attributes (id, name, type, required, default_value, options, merchant_id, enabled, searchable, deleted, code) FROM stdin;
\.


--
-- Data for Name: brands; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.brands (id, name, thumbnail, description, deleted, enabled, sort_order, created_at, updated_at, type, machine_types) FROM stdin;
\.


--
-- Data for Name: carts; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.carts (id, product_id, user_id, product_count, merchant_id, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: cash_request; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.cash_request (id, amount, debit_card_info, status, reply, note) FROM stdin;
\.


--
-- Data for Name: categories; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.categories (id, parent_id, name, sort_order, thumbnail, deleted, enabled, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: comments; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.comments (id, user_id, product_id, rank, order_id, content, is_show, created_at, user_ip, reply, merchant_id) FROM stdin;
\.


--
-- Data for Name: configs; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.configs (id, code, value, sort_order, name, description, deleted) FROM stdin;
\.


--
-- Data for Name: debit_card_info; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.debit_card_info (user_id, debit_card_info, is_checked, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: favorites; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.favorites (id, user_id, created_at, type, object_id) FROM stdin;
\.


--
-- Data for Name: goose_db_version; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.goose_db_version (id, version_id, is_applied, tstamp) FROM stdin;
\.


--
-- Data for Name: merchant_profiles; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.merchant_profiles (user_id, social_id, rep_name, company_name, company_address, delivery_address, license_image, company_image, created_at, updated_at) FROM stdin;
\.


--
-- Data for Name: operation_logs; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.operation_logs (id, user_id, action, created_at) FROM stdin;
\.


--
-- Data for Name: order_products; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.order_products (id, order_id, product_id, product_name, product_number, product_price, product_image) FROM stdin;
\.


--
-- Data for Name: orders; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.orders (id, user_id, order_status, shipping_status, pay_status, postscript, shipping_id, shipping_name, pay_id, pay_name, how_oos, inv_payee, inv_type, inv_content, amount, shipping_fee, pay_fee, money_paid, order_amount, created_at, confirm_time, pay_time, shipping_time, tax, parent_id, merchant_id, address, note) FROM stdin;
\.


--
-- Data for Name: payments; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.payments (id, pay_name, pay_code, pay_fee, pay_desc, sort_order, enabled, is_cod, is_online, deleted) FROM stdin;
\.


--
-- Data for Name: product_images; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.product_images (id, product_id, small_image, medium_image, big_image) FROM stdin;
\.


--
-- Data for Name: products; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.products (id, name, is_sale, attribute_set_id, batch_price, second_price, category_id, related_ids, content, brand_id, deleted, created_at, updated_at, tags, attrs, recommended, sales_volume, machine_types) FROM stdin;
\.


--
-- Data for Name: regions; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.regions (id, parent_id, name, sort_order) FROM stdin;
\.


--
-- Data for Name: rel_merchants_products; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.rel_merchants_products (product_id, merchant_id, stock, retail_price) FROM stdin;
\.


--
-- Data for Name: users; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.users (id, name, mobile, password, type, email, created_at, updated_at, status, last_ip, last_login_time) FROM stdin;
\.


--
-- Data for Name: wecharts; Type: TABLE DATA; Schema: public; Owner: postgres
--

COPY public.wecharts (openid, session_key, nick_name, gender, language, city, province, country, avatar_url, user_id) FROM stdin;
\.


--
-- Name: goose_db_version_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.goose_db_version_id_seq', 1, false);


--
-- Name: regions_id_seq; Type: SEQUENCE SET; Schema: public; Owner: postgres
--

SELECT pg_catalog.setval('public.regions_id_seq', 1, false);


--
-- Name: addresses addresses_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.addresses
    ADD CONSTRAINT addresses_pkey PRIMARY KEY (id);


--
-- Name: attributes attr_code; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.attributes
    ADD CONSTRAINT attr_code UNIQUE (code);


--
-- Name: attribute_sets attribute_sets_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.attribute_sets
    ADD CONSTRAINT attribute_sets_pkey PRIMARY KEY (id);


--
-- Name: attributes attributes_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.attributes
    ADD CONSTRAINT attributes_pkey PRIMARY KEY (id);


--
-- Name: brands brands_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.brands
    ADD CONSTRAINT brands_pkey PRIMARY KEY (id);


--
-- Name: carts cart_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.carts
    ADD CONSTRAINT cart_pkey PRIMARY KEY (id);


--
-- Name: cash_request cash_request_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.cash_request
    ADD CONSTRAINT cash_request_pkey PRIMARY KEY (id);


--
-- Name: categories categories_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT categories_pkey PRIMARY KEY (id);


--
-- Name: comments comments_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.comments
    ADD CONSTRAINT comments_pkey PRIMARY KEY (id);


--
-- Name: configs config_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.configs
    ADD CONSTRAINT config_pkey PRIMARY KEY (id);


--
-- Name: debit_card_info debit_card_info_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.debit_card_info
    ADD CONSTRAINT debit_card_info_pkey PRIMARY KEY (user_id);


--
-- Name: favorites favorites_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.favorites
    ADD CONSTRAINT favorites_pkey PRIMARY KEY (id);


--
-- Name: goose_db_version goose_db_version_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.goose_db_version
    ADD CONSTRAINT goose_db_version_pkey PRIMARY KEY (id);


--
-- Name: merchant_profiles merchant_profiles_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.merchant_profiles
    ADD CONSTRAINT merchant_profiles_pkey PRIMARY KEY (user_id);


--
-- Name: operation_logs operation_logs_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.operation_logs
    ADD CONSTRAINT operation_logs_pkey PRIMARY KEY (id);


--
-- Name: order_products order_products_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.order_products
    ADD CONSTRAINT order_products_pkey PRIMARY KEY (id);


--
-- Name: orders orders_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.orders
    ADD CONSTRAINT orders_pkey PRIMARY KEY (id);


--
-- Name: payments payment_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.payments
    ADD CONSTRAINT payment_pkey PRIMARY KEY (id);


--
-- Name: product_images product_images_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.product_images
    ADD CONSTRAINT product_images_pkey PRIMARY KEY (id);


--
-- Name: products products_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_pkey PRIMARY KEY (id);


--
-- Name: regions regions_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.regions
    ADD CONSTRAINT regions_pkey PRIMARY KEY (id);


--
-- Name: rel_merchants_products rel_merchants_products_user_id_product_id_pk; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.rel_merchants_products
    ADD CONSTRAINT rel_merchants_products_user_id_product_id_pk PRIMARY KEY (merchant_id, product_id);


--
-- Name: users users_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (id);


--
-- Name: wecharts wecharts_pkey; Type: CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.wecharts
    ADD CONSTRAINT wecharts_pkey PRIMARY KEY (openid);


--
-- Name: products_category_id_index; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX products_category_id_index ON public.products USING btree (category_id);


--
-- Name: products_machine_types_index; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX products_machine_types_index ON public.products USING btree (machine_types);


--
-- Name: products_sales_volume_index; Type: INDEX; Schema: public; Owner: postgres
--

CREATE INDEX products_sales_volume_index ON public.products USING btree (sales_volume DESC);


--
-- Name: addresses update_time; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER update_time BEFORE UPDATE ON public.addresses FOR EACH ROW EXECUTE PROCEDURE public.auto_update_timestamp();


--
-- Name: brands update_time; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER update_time BEFORE UPDATE ON public.brands FOR EACH ROW EXECUTE PROCEDURE public.auto_update_timestamp();


--
-- Name: carts update_time; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER update_time BEFORE UPDATE ON public.carts FOR EACH ROW EXECUTE PROCEDURE public.auto_update_timestamp();


--
-- Name: categories update_time; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER update_time BEFORE UPDATE ON public.categories FOR EACH ROW EXECUTE PROCEDURE public.auto_update_timestamp();


--
-- Name: merchant_profiles update_time; Type: TRIGGER; Schema: public; Owner: postgres
--

CREATE TRIGGER update_time BEFORE UPDATE ON public.merchant_profiles FOR EACH ROW EXECUTE PROCEDURE public.auto_update_timestamp();


--
-- Name: addresses addresses_users_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.addresses
    ADD CONSTRAINT addresses_users_id_fk FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: carts carts_merchant_profiles_user_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.carts
    ADD CONSTRAINT carts_merchant_profiles_user_id_fk FOREIGN KEY (merchant_id) REFERENCES public.merchant_profiles(user_id);


--
-- Name: carts carts_products_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.carts
    ADD CONSTRAINT carts_products_id_fk FOREIGN KEY (product_id) REFERENCES public.products(id);


--
-- Name: carts carts_users_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.carts
    ADD CONSTRAINT carts_users_id_fk FOREIGN KEY (user_id) REFERENCES public.users(id);


--
-- Name: products products_brands_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_brands_id_fk FOREIGN KEY (brand_id) REFERENCES public.brands(id);


--
-- Name: products products_categories_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.products
    ADD CONSTRAINT products_categories_id_fk FOREIGN KEY (category_id) REFERENCES public.categories(id);


--
-- Name: rel_merchants_products rel_merchants_products_products_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.rel_merchants_products
    ADD CONSTRAINT rel_merchants_products_products_id_fk FOREIGN KEY (product_id) REFERENCES public.products(id);


--
-- Name: rel_merchants_products rel_merchants_products_users_id_fk; Type: FK CONSTRAINT; Schema: public; Owner: postgres
--

ALTER TABLE ONLY public.rel_merchants_products
    ADD CONSTRAINT rel_merchants_products_users_id_fk FOREIGN KEY (merchant_id) REFERENCES public.users(id);


--
-- PostgreSQL database dump complete
--

