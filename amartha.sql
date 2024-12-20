PGDMP  	    #            
    |            amartha    16.0    16.0 *    :           0    0    ENCODING    ENCODING        SET client_encoding = 'UTF8';
                      false            ;           0    0 
   STDSTRINGS 
   STDSTRINGS     (   SET standard_conforming_strings = 'on';
                      false            <           0    0 
   SEARCHPATH 
   SEARCHPATH     8   SELECT pg_catalog.set_config('search_path', '', false);
                      false            =           1262    17772    amartha    DATABASE     i   CREATE DATABASE amartha WITH TEMPLATE = template0 ENCODING = 'UTF8' LOCALE_PROVIDER = libc LOCALE = 'C';
    DROP DATABASE amartha;
                postgres    false                        2615    2200    public    SCHEMA        CREATE SCHEMA public;
    DROP SCHEMA public;
                pg_database_owner    false            >           0    0    SCHEMA public    COMMENT     6   COMMENT ON SCHEMA public IS 'standard public schema';
                   pg_database_owner    false    4            `           1247    17844    loan_status    TYPE     l   CREATE TYPE public.loan_status AS ENUM (
    'proposed',
    'approved',
    'invested',
    'disbursed'
);
    DROP TYPE public.loan_status;
       public          postgres    false    4            �            1259    17824    admin    TABLE     (  CREATE TABLE public.admin (
    id bigint NOT NULL,
    name character varying NOT NULL,
    email character varying NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone,
    deleted_at timestamp without time zone
);
    DROP TABLE public.admin;
       public         heap    postgres    false    4            �            1259    17823    admin_id_seq    SEQUENCE     u   CREATE SEQUENCE public.admin_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
 #   DROP SEQUENCE public.admin_id_seq;
       public          postgres    false    4    224            ?           0    0    admin_id_seq    SEQUENCE OWNED BY     =   ALTER SEQUENCE public.admin_id_seq OWNED BY public.admin.id;
          public          postgres    false    223            �            1259    17814    borrower    TABLE     y  CREATE TABLE public.borrower (
    id bigint NOT NULL,
    name character varying NOT NULL,
    email character varying NOT NULL,
    phone character varying NOT NULL,
    address character varying NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone,
    deleted_at timestamp without time zone
);
    DROP TABLE public.borrower;
       public         heap    postgres    false    4            �            1259    17813    borrower_id_seq    SEQUENCE     x   CREATE SEQUENCE public.borrower_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
 &   DROP SEQUENCE public.borrower_id_seq;
       public          postgres    false    222    4            @           0    0    borrower_id_seq    SEQUENCE OWNED BY     C   ALTER SEQUENCE public.borrower_id_seq OWNED BY public.borrower.id;
          public          postgres    false    221            �            1259    17804    investor    TABLE     Q  CREATE TABLE public.investor (
    id bigint NOT NULL,
    name character varying NOT NULL,
    email character varying NOT NULL,
    phone character varying NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone,
    deleted_at timestamp without time zone
);
    DROP TABLE public.investor;
       public         heap    postgres    false    4            �            1259    17803    investor_id_seq    SEQUENCE     x   CREATE SEQUENCE public.investor_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
 &   DROP SEQUENCE public.investor_id_seq;
       public          postgres    false    220    4            A           0    0    investor_id_seq    SEQUENCE OWNED BY     C   ALTER SEQUENCE public.investor_id_seq OWNED BY public.investor.id;
          public          postgres    false    219            �            1259    17787    loan    TABLE     �  CREATE TABLE public.loan (
    id bigint NOT NULL,
    borrower_id bigint NOT NULL,
    principal_amount bigint NOT NULL,
    status public.loan_status DEFAULT 'proposed'::public.loan_status NOT NULL,
    rate numeric(5,2) NOT NULL,
    approved_by bigint,
    approved_at timestamp without time zone,
    agreement_letter character varying,
    field_validator character varying,
    roi numeric(5,2) NOT NULL,
    disbursed_by bigint,
    disbursed_at timestamp without time zone,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone,
    deleted_at timestamp without time zone
);
    DROP TABLE public.loan;
       public         heap    postgres    false    864    864    4            �            1259    17796    loan_detail    TABLE     F  CREATE TABLE public.loan_detail (
    id bigint NOT NULL,
    loan_id bigint NOT NULL,
    investor_id bigint NOT NULL,
    invested_amount bigint NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone,
    deleted_at timestamp without time zone
);
    DROP TABLE public.loan_detail;
       public         heap    postgres    false    4            �            1259    17795    loan_detail_id_seq    SEQUENCE     {   CREATE SEQUENCE public.loan_detail_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
 )   DROP SEQUENCE public.loan_detail_id_seq;
       public          postgres    false    218    4            B           0    0    loan_detail_id_seq    SEQUENCE OWNED BY     I   ALTER SEQUENCE public.loan_detail_id_seq OWNED BY public.loan_detail.id;
          public          postgres    false    217            �            1259    17786    loan_id_seq    SEQUENCE     t   CREATE SEQUENCE public.loan_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;
 "   DROP SEQUENCE public.loan_id_seq;
       public          postgres    false    4    216            C           0    0    loan_id_seq    SEQUENCE OWNED BY     ;   ALTER SEQUENCE public.loan_id_seq OWNED BY public.loan.id;
          public          postgres    false    215            �           2604    17827    admin id    DEFAULT     d   ALTER TABLE ONLY public.admin ALTER COLUMN id SET DEFAULT nextval('public.admin_id_seq'::regclass);
 7   ALTER TABLE public.admin ALTER COLUMN id DROP DEFAULT;
       public          postgres    false    223    224    224            �           2604    17817    borrower id    DEFAULT     j   ALTER TABLE ONLY public.borrower ALTER COLUMN id SET DEFAULT nextval('public.borrower_id_seq'::regclass);
 :   ALTER TABLE public.borrower ALTER COLUMN id DROP DEFAULT;
       public          postgres    false    221    222    222            �           2604    17807    investor id    DEFAULT     j   ALTER TABLE ONLY public.investor ALTER COLUMN id SET DEFAULT nextval('public.investor_id_seq'::regclass);
 :   ALTER TABLE public.investor ALTER COLUMN id DROP DEFAULT;
       public          postgres    false    220    219    220            �           2604    17790    loan id    DEFAULT     b   ALTER TABLE ONLY public.loan ALTER COLUMN id SET DEFAULT nextval('public.loan_id_seq'::regclass);
 6   ALTER TABLE public.loan ALTER COLUMN id DROP DEFAULT;
       public          postgres    false    216    215    216            �           2604    17799    loan_detail id    DEFAULT     p   ALTER TABLE ONLY public.loan_detail ALTER COLUMN id SET DEFAULT nextval('public.loan_detail_id_seq'::regclass);
 =   ALTER TABLE public.loan_detail ALTER COLUMN id DROP DEFAULT;
       public          postgres    false    217    218    218            7          0    17824    admin 
   TABLE DATA           T   COPY public.admin (id, name, email, created_at, updated_at, deleted_at) FROM stdin;
    public          postgres    false    224   +0       5          0    17814    borrower 
   TABLE DATA           g   COPY public.borrower (id, name, email, phone, address, created_at, updated_at, deleted_at) FROM stdin;
    public          postgres    false    222   |0       3          0    17804    investor 
   TABLE DATA           ^   COPY public.investor (id, name, email, phone, created_at, updated_at, deleted_at) FROM stdin;
    public          postgres    false    220   �0       /          0    17787    loan 
   TABLE DATA           �   COPY public.loan (id, borrower_id, principal_amount, status, rate, approved_by, approved_at, agreement_letter, field_validator, roi, disbursed_by, disbursed_at, created_at, updated_at, deleted_at) FROM stdin;
    public          postgres    false    216   X1       1          0    17796    loan_detail 
   TABLE DATA           t   COPY public.loan_detail (id, loan_id, investor_id, invested_amount, created_at, updated_at, deleted_at) FROM stdin;
    public          postgres    false    218   u1       D           0    0    admin_id_seq    SEQUENCE SET     :   SELECT pg_catalog.setval('public.admin_id_seq', 1, true);
          public          postgres    false    223            E           0    0    borrower_id_seq    SEQUENCE SET     =   SELECT pg_catalog.setval('public.borrower_id_seq', 1, true);
          public          postgres    false    221            F           0    0    investor_id_seq    SEQUENCE SET     =   SELECT pg_catalog.setval('public.investor_id_seq', 2, true);
          public          postgres    false    219            G           0    0    loan_detail_id_seq    SEQUENCE SET     A   SELECT pg_catalog.setval('public.loan_detail_id_seq', 1, false);
          public          postgres    false    217            H           0    0    loan_id_seq    SEQUENCE SET     :   SELECT pg_catalog.setval('public.loan_id_seq', 1, false);
          public          postgres    false    215            �           2606    17832    admin admin_pkey 
   CONSTRAINT     N   ALTER TABLE ONLY public.admin
    ADD CONSTRAINT admin_pkey PRIMARY KEY (id);
 :   ALTER TABLE ONLY public.admin DROP CONSTRAINT admin_pkey;
       public            postgres    false    224            �           2606    17822    borrower borrower_pkey 
   CONSTRAINT     T   ALTER TABLE ONLY public.borrower
    ADD CONSTRAINT borrower_pkey PRIMARY KEY (id);
 @   ALTER TABLE ONLY public.borrower DROP CONSTRAINT borrower_pkey;
       public            postgres    false    222            �           2606    17812    investor investor_pkey 
   CONSTRAINT     T   ALTER TABLE ONLY public.investor
    ADD CONSTRAINT investor_pkey PRIMARY KEY (id);
 @   ALTER TABLE ONLY public.investor DROP CONSTRAINT investor_pkey;
       public            postgres    false    220            �           2606    17802    loan_detail loan_detail_pkey 
   CONSTRAINT     Z   ALTER TABLE ONLY public.loan_detail
    ADD CONSTRAINT loan_detail_pkey PRIMARY KEY (id);
 F   ALTER TABLE ONLY public.loan_detail DROP CONSTRAINT loan_detail_pkey;
       public            postgres    false    218            �           2606    17794    loan loan_pkey 
   CONSTRAINT     L   ALTER TABLE ONLY public.loan
    ADD CONSTRAINT loan_pkey PRIMARY KEY (id);
 8   ALTER TABLE ONLY public.loan DROP CONSTRAINT loan_pkey;
       public            postgres    false    216            7   A   x�3�N�-.��,S鹉�9z����FF&������
�fV�@d�gfn`fn��D\1z\\\ b,O      5   V   x�3�IMI��,�鹉�9z��������f��%�y�E@�����D��P��R��������@�(hh��D\1z\\\ ��      3   f   x�3���+K-.�/R0�3�s3s���s9,-�����9��Ltu���L���L��-,9c����a�iD�q&fz���&&fP�b���� a{)�      /      x������ � �      1      x������ � �     