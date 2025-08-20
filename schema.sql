--
-- PostgreSQL database dump
--

-- Dumped from database version 17.5
-- Dumped by pg_dump version 17.5

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET transaction_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SELECT pg_catalog.set_config('search_path', '', false);
SET check_function_bodies = false;
SET xmloption = content;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: email; Type: DOMAIN; Schema: public; Owner: -
--

CREATE DOMAIN public.email AS character varying(200)
	CONSTRAINT email_check CHECK (((VALUE)::text ~* '^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$'::text));


--
-- Name: phone_number; Type: DOMAIN; Schema: public; Owner: -
--

CREATE DOMAIN public.phone_number AS character varying(100)
	CONSTRAINT phone_number_check CHECK (((VALUE)::text ~ '^09[0-9]{9}$'::text));


--
-- Name: status; Type: DOMAIN; Schema: public; Owner: -
--

CREATE DOMAIN public.status AS character varying(20)
	CONSTRAINT status_check CHECK (((VALUE)::text = ANY ((ARRAY['active'::character varying, 'suspended'::character varying])::text[])));


--
-- Name: prevent_insert_loan_as_returend(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.prevent_insert_loan_as_returend() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
begin
	if new.status = 'returned' then
		raise exception 'NON_BORROWED_STATUS';
	end if;
end;
$$;


--
-- Name: set_available_copies_to_total(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.set_available_copies_to_total() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
begin
	NEW.available_copies = NEW.total_copies;
	return NEW;
end;
$$;


--
-- Name: update_available_copies(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.update_available_copies() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
declare
	new_available_copies INT;
	book_total_copies INT;
begin
	-- lock book row for updating
	PERFORM 1 FROM books where id = COALESCE(new.book_id,old.book_id);

	-- get book's total copies
	select total_copies into book_total_copies 
	from books
	where id = COALESCE(new.book_id,old.book_id);

	-- calculate active loans and new avail able copies
	select book_total_copies - count(*) 
	into new_available_copies
	from loans
	where book_id = COALESCE(new.book_id,old.book_id) and status = 'borrowed';

	if new_available_copies < 0  then
		raise exception 'OVER_BORROWING';
	end if;	
	
	-- update book
	update books set available_copies = new_available_copies
	where id = COALESCE(new.book_id,old.book_id);
		
	return NEW;
end;
$$;


--
-- Name: update_available_copies_on_totalcopy_update(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.update_available_copies_on_totalcopy_update() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
declare
	active_loans int;
begin
	-- lock
	perform 1 from books where id = NEW.id for update;
	-- count active loans
	select count(*) into active_loans from loans
	where book_id = NEW.id and status = 'borrowed';
	
	if NEW.total_copies < active_loans then
		raise exception 'INVALID_TOTAL_COPIES';
	end if;

	NEW.available_copies = NEW.total_copies - active_loans;
	return NEW;
end; 
$$;


--
-- Name: update_update_fields(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.update_update_fields() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
begin
	NEW.updated_at := now();
	NEW.version := OLD.version + 1;
	return NEW;
end;
$$;


--
-- Name: validate_available_copies_update(); Type: FUNCTION; Schema: public; Owner: -
--

CREATE FUNCTION public.validate_available_copies_update() RETURNS trigger
    LANGUAGE plpgsql
    AS $$
declare
	active_loans int;
begin
	SELECT count(*) into active_loans from loans
	where book_id = new.id and status = 'borrowed';
	
	if new.available_copies != (new.total_copies - active_loans) then
		raise exception 'INVALID_AVAILABLE_COPIES';
	end if;

	return NEW;

end;
$$;


--
-- Name: active_members_view; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.active_members_view AS
SELECT
    NULL::integer AS id,
    NULL::text AS full_name,
    NULL::bigint AS loan_count;


SET default_tablespace = '';

SET default_table_access_method = heap;

--
-- Name: activity_logs; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.activity_logs (
    id integer NOT NULL,
    activity_type character varying(50) NOT NULL,
    actor_id integer,
    actor_type character varying(50) NOT NULL,
    description text NOT NULL,
    entity_id integer,
    entity_type character varying(50),
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    CONSTRAINT activity_log_actor_type_check CHECK (((actor_type)::text = ANY ((ARRAY['staff'::character varying, 'member'::character varying])::text[]))),
    CONSTRAINT activity_log_entity_type_check CHECK (((entity_type)::text = ANY ((ARRAY['author'::character varying, 'book'::character varying, 'category'::character varying, 'loan'::character varying, 'member'::character varying, 'staff'::character varying])::text[])))
);


--
-- Name: activity_log_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.activity_log_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: activity_log_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.activity_log_id_seq OWNED BY public.activity_logs.id;


--
-- Name: authors; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.authors (
    id integer NOT NULL,
    full_name text NOT NULL,
    nationality text NOT NULL,
    bio text,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    version integer DEFAULT 1 NOT NULL
);


--
-- Name: authors_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.authors_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: authors_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.authors_id_seq OWNED BY public.authors.id;


--
-- Name: books; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.books (
    id integer NOT NULL,
    title text NOT NULL,
    isbn character varying(20) NOT NULL,
    author_id integer NOT NULL,
    total_copies integer NOT NULL,
    available_copies integer NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    version integer DEFAULT 1 NOT NULL,
    publisher text,
    language text,
    summary text,
    translator text,
    cover_image_url text,
    CONSTRAINT books_available_copies_check CHECK ((available_copies >= 0)),
    CONSTRAINT books_total_copies_check CHECK ((total_copies >= 0))
);


--
-- Name: books_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.books_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: books_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.books_id_seq OWNED BY public.books.id;


--
-- Name: categories; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.categories (
    id integer NOT NULL,
    slug character varying(100) NOT NULL,
    name character varying(100) NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    version integer DEFAULT 1 NOT NULL
);


--
-- Name: category_books; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.category_books (
    book_id integer NOT NULL,
    category_id integer NOT NULL
);


--
-- Name: category_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.category_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: category_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.category_id_seq OWNED BY public.categories.id;


--
-- Name: loans; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.loans (
    id integer NOT NULL,
    book_id integer NOT NULL,
    member_id integer NOT NULL,
    borrow_date timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    due_date timestamp without time zone NOT NULL,
    return_date timestamp without time zone,
    status text DEFAULT 'borrowed'::text,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    version integer DEFAULT 1 NOT NULL,
    CONSTRAINT loans_status_check CHECK ((status = ANY (ARRAY['borrowed'::text, 'returned'::text])))
);


--
-- Name: loans_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.loans_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: loans_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.loans_id_seq OWNED BY public.loans.id;


--
-- Name: low_stock_high_demand_books_view; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.low_stock_high_demand_books_view AS
SELECT
    NULL::integer AS id,
    NULL::text AS title,
    NULL::integer AS total_copies,
    NULL::integer AS available_copies,
    NULL::bigint AS loan_count;


--
-- Name: members; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.members (
    id integer NOT NULL,
    full_name text NOT NULL,
    email public.email NOT NULL,
    phone_number public.phone_number NOT NULL,
    joined_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP,
    status public.status DEFAULT 'active'::text,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    version integer DEFAULT 1 NOT NULL,
    national_id character varying(20) NOT NULL,
    CONSTRAINT members_status_check CHECK (((status)::text = ANY (ARRAY['active'::text, 'suspended'::text])))
);


--
-- Name: members_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.members_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: members_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.members_id_seq OWNED BY public.members.id;


--
-- Name: popular_authors_view; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.popular_authors_view AS
SELECT
    NULL::integer AS id,
    NULL::text AS full_name,
    NULL::bigint AS total_loans;


--
-- Name: popular_books_view; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.popular_books_view AS
SELECT
    NULL::integer AS id,
    NULL::text AS title,
    NULL::bigint AS loan_count;


--
-- Name: popular_categories_view; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.popular_categories_view AS
 SELECT cg.name,
    cg.slug,
    count(l.id) AS loan_count
   FROM (((public.categories cg
     JOIN public.category_books cb ON ((cb.category_id = cg.id)))
     JOIN public.books b ON ((cb.book_id = b.id)))
     JOIN public.loans l ON ((l.book_id = b.id)))
  GROUP BY cg.name, cg.slug
  ORDER BY (count(l.id)) DESC;


--
-- Name: staff; Type: TABLE; Schema: public; Owner: -
--

CREATE TABLE public.staff (
    id integer NOT NULL,
    full_name character varying(100) NOT NULL,
    phone_number public.phone_number NOT NULL,
    email public.email NOT NULL,
    role character varying(50) DEFAULT 'librarian'::character varying NOT NULL,
    password_hash text NOT NULL,
    created_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at timestamp without time zone DEFAULT CURRENT_TIMESTAMP NOT NULL,
    last_login timestamp without time zone,
    version integer DEFAULT 1 NOT NULL,
    status public.status DEFAULT 'active'::character varying NOT NULL,
    CONSTRAINT staff_role_check CHECK (((role)::text = ANY ((ARRAY['admin'::character varying, 'librarian'::character varying])::text[])))
);


--
-- Name: staff_id_seq; Type: SEQUENCE; Schema: public; Owner: -
--

CREATE SEQUENCE public.staff_id_seq
    AS integer
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


--
-- Name: staff_id_seq; Type: SEQUENCE OWNED BY; Schema: public; Owner: -
--

ALTER SEQUENCE public.staff_id_seq OWNED BY public.staff.id;


--
-- Name: upcoming_due_loans; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.upcoming_due_loans AS
 SELECT l.id,
    b.title AS book_title,
    m.full_name AS member_name,
    l.due_date
   FROM ((public.loans l
     JOIN public.books b ON ((b.id = l.book_id)))
     JOIN public.members m ON ((m.id = l.member_id)))
  WHERE ((l.return_date IS NULL) AND (l.due_date >= CURRENT_DATE) AND (l.due_date <= (CURRENT_DATE + '2 days'::interval)))
  ORDER BY l.due_date;


--
-- Name: upcoming_due_loans_view; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.upcoming_due_loans_view AS
 SELECT l.id,
    b.title AS book_title,
    m.full_name AS member_name,
    l.due_date
   FROM ((public.loans l
     JOIN public.books b ON ((b.id = l.book_id)))
     JOIN public.members m ON ((m.id = l.member_id)))
  WHERE ((l.return_date IS NULL) AND (l.due_date >= CURRENT_DATE) AND (l.due_date <= (CURRENT_DATE + '2 days'::interval)))
  ORDER BY l.due_date;


--
-- Name: upcoming_loans_view; Type: VIEW; Schema: public; Owner: -
--

CREATE VIEW public.upcoming_loans_view AS
 SELECT l.id,
    b.title AS book_title,
    m.full_name AS member_name,
    l.due_date
   FROM ((public.loans l
     JOIN public.books b ON ((b.id = l.book_id)))
     JOIN public.members m ON ((m.id = l.member_id)))
  WHERE ((l.return_date IS NULL) AND (l.due_date >= CURRENT_DATE) AND (l.due_date <= (CURRENT_DATE + '2 days'::interval)))
  ORDER BY l.due_date;


--
-- Name: activity_logs id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.activity_logs ALTER COLUMN id SET DEFAULT nextval('public.activity_log_id_seq'::regclass);


--
-- Name: authors id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.authors ALTER COLUMN id SET DEFAULT nextval('public.authors_id_seq'::regclass);


--
-- Name: books id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.books ALTER COLUMN id SET DEFAULT nextval('public.books_id_seq'::regclass);


--
-- Name: categories id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.categories ALTER COLUMN id SET DEFAULT nextval('public.category_id_seq'::regclass);


--
-- Name: loans id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.loans ALTER COLUMN id SET DEFAULT nextval('public.loans_id_seq'::regclass);


--
-- Name: members id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.members ALTER COLUMN id SET DEFAULT nextval('public.members_id_seq'::regclass);


--
-- Name: staff id; Type: DEFAULT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.staff ALTER COLUMN id SET DEFAULT nextval('public.staff_id_seq'::regclass);


--
-- Data for Name: activity_logs; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.activity_logs (id, activity_type, actor_id, actor_type, description, entity_id, entity_type, created_at) FROM stdin;
\.


--
-- Data for Name: authors; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.authors (id, full_name, nationality, bio, created_at, updated_at, version) FROM stdin;
14	گابریل گارسیا مارکز	کلمبیا	گابریل خوزه گارسیا مارکِز (اسپانیایی: Gabriel José García Márquez‎؛ ۶ مارس ۱۹۲۷ – ۱۷ آوریل ۲۰۱۴) رمان‌نویس، نویسنده، روزنامه‌نگار، ناشر و فعال سیاسی کلمبیایی بود. او بین مردم کشورهای آمریکای لاتین با نام گابو یا گابیتو (برای محبوبیت) مشهور بود و پس از درگیری با رئیس دولت کلمبیا و تحت تعقیب قرار گرفتن، در مکزیک زندگی می‌کرد. مارکز برنده جایزه نوبل ادبیات در سال ۱۹۸۲ بود. او را بیش از سایر آثارش به خاطر رمان صد سال تنهایی چاپ ۱۹۶۷ می‌شناسند که یکی از پرفروش‌ترین کتاب‌های جهان در سبک رئالیسم جادویی است.	2025-08-19 12:56:31.794049	2025-08-19 12:56:31.794049	1
11	جورج اورول	انگلیس	اریک آرتور بلر (به انگلیسی: Eric Arthur Blair) با نامِ مستعار جُرج اوروِل (به انگلیسی: George Orwell) (۲۵ ژوئن ۱۹۰۳ – ۲۱ ژانویه ۱۹۵۰) داستان‌نویس، روزنامه‌نگار، منتقدِ ادبی و شاعر انگلیسی بود. او بیشتر برای دو رمان سرشناس و پرفروش مزرعه حیوانات که در ۱۹۴۵ منتشر شد و در اواخر دههٔ ۱۹۵۰ به شهرت رسید و نیز رمان ۱۹۸۴ شناخته می‌شود. این دو کتاب بر روی هم بیش از هر دو کتابِ دیگری از یک نویسندهٔ قرن بیستمی، فروش داشته‌اند.[۲] او همچنین با نقدهای پرشماری که بر کتاب‌ها نوشت، بهترین وقایع‌نگار فرهنگ و ادب انگلیسی قرن شناخته می‌شود.\r\n\r\n	2025-08-19 12:25:17.347666	2025-08-19 08:55:30.243343	2
12	هارپر لی	آمریکا	نل هارپر لی (به انگلیسی: Nelle Harper Lee) (زاده ۲۸ آوریل ۱۹۲۶ – درگذشته ۱۹ فوریه ۲۰۱۶) رمان‌نویس اهل ایالات متحده آمریکا بود که در سال ۱۹۶۰ به خاطر نگارش کتاب کشتن مرغ مقلد برنده جایزه پولیتزر برای رمان شد.	2025-08-19 12:45:58.75969	2025-08-19 12:45:58.75969	1
13	اف. اسکات فیتزجرالد	آمریکا	فرانسیس اسکات کی فیتزجرالد (انگلیسی: Francis Scott Key Fitzgerald؛ ۲۴ سپتامبر ۱۸۹۶ – ۲۱ دسامبر ۱۹۴۰) نویسنده آمریکایی رمان و داستان‌های کوتاه بود. آثار فیتزجرالد نمایانگر عصر جاز در آمریکا است. او به عنوان یکی از نویسندگان بزرگ سده بیستم میلادی شناخته می‌شود. تاثیرگذارترین اثر او رمان گتسبی بزرگ است که اولین بار در سال ۱۹۲۵ منتشر شد.[۱] 	2025-08-19 12:52:37.56286	2025-08-19 12:52:37.56286	1
15	فیودور داستایفسکی	روسیه	فیودور میخایلوویچ داستایِفسکی یا فیودور میخایلوویچ دوستویِوسکی (روسی: Фёдор Михайлович Достоевский[الف]؛ IPA: [ˈfʲɵdər mʲɪˈxajləvʲɪdʑ dəstɐˈjɛfskʲɪj] (شنیدنⓘ) ؛ زادهٔ ۱۱ نوامبر ۱۸۲۱ – درگذشتهٔ ۹ فوریهٔ ۱۸۸۱[۱][ب]) نویسندهٔ مشهور و تأثیرگذار روسی بود. داستایفسکی عمق و فراز جامعهٔ روسیه را تجربه کرد. شخصیت‌های او در همهٔ رمان‌هایش با مشکلات روان‌شناسانه و عاطفی درگیر هستند اما مهم‌تر آنکه کتاب‌هایش از آموزه‌های ایدئولوژیک زمان خود الهام می‌گرفتند. رمان‌های مشهور او جنایت و مکافات (۱۸۶۶)، ابله (۱۸۶۹)، جن‌زدگان (۱۸۷۲) و برادران کارامازوف (۱۸۸۰) هستند.	2025-08-19 13:00:35.894279	2025-08-19 13:00:35.894279	1
16	لئو تولستوی	روسیه	لف نیکلایِویچ تولستوی (به روسی: Лев Никола́евич Толсто́й)  (۹ سپتامبر ۱۸۲۸–۲۰ نوامبر ۱۹۱۰) فیلسوف، عارف و نویسندۀ روس بود. او را از بزرگ‌ترین رمان‌نویسان جهان می‌دانند.[۱] او بارها نامزد دریافت جایزۀ نوبل ادبیات و جایزۀ صلح نوبل شد؛ ولی هرگز به آنها دست نیافت.\r\n\r\nرمان‌های جنگ و صلح و آنا کارنینا که همواره در بین بهترین رمان‌های جهان هستند، آثار تولستوی‌اند. او در روسیه هواداران بسیاری دارد و سکۀ طلای یادبود برای بزرگداشت وی ضرب شده‌ است. تولستوی در زمان زندگی خود در جهان سرشناس، ولی ساده‌زیست بود.\r\n\r\n	2025-08-19 13:04:26.529379	2025-08-19 13:04:26.529379	1
17	میگل د سروانتس	اسپانیا	میگل د سِروانتِس ساآودرا (به انگلیسی: Miguel de Cervantes Saavedra)، (اسپانیایی: [miˈɣel de θeɾˈβantes saaˈβeðɾa]؛ زاده ۲۹ سپتامبر ۱۵۴۷ (با فرض) – درگذشته ۲۲ آوریل ۱۶۱۶ NS)[۲] رمان‌نویس، شاعر، نقاش و نمایشنامه‌نویس نامدار اسپانیایی بود.\r\n\r\nرمان مشهور دن کیشوت - که از پایه‌های ادبیات کلاسیک اروپا به‌شمار می‌آید و بسیاری از منتقدان از آن به عنوان نخستین رمان مدرن و یکی از بهترین آثار ادبی جهان یاد می‌کنند - اثر اوست.[۳] به وی لقب شاهزادهٔ نبوغ داده‌اند. برخی رمان دن کیشوت را تاثیرگذارترین رمان ژانر تخیلی (Fiction) می‌دانند.[۴]\r\n\r\nجالب این است که سروانتس کتاب دن کیشوت معروف که تمام اروپا و جهان را در زمان خود معطوف نموده بود را در زندان و در دوران گذراندن زمان زندانی خود نوشته بود.	2025-08-19 13:08:06.951108	2025-08-19 13:08:06.951108	1
18	گوستاو فلوبر	فرانسه	گوستاو فلوبر (به فرانسوی: Gustave Flaubert) (زاده ۱۲ دسامبر ۱۸۲۱ – درگذشته ۸ مه ۱۸۸۰) از نویسندگان تأثیرگذار قرن نوزدهم فرانسه بود که اغلب جزو بزرگترین رمان نویسان ادبیات غرب شمرده می‌شود. نوع نگارش واقع‌گرایانهٔ فلوبر، ادبیات بسیار غنی و تحلیل‌های روان‌شناختی عمیق او از جمله خصوصیات آثار وی است که الهام‌بخش نویسندگانی چون گی دو موپاسان، امیل زولا و آلفونس دوده بوده‌است. او خود تأثیرگرفته از سبک و موضوعات بالزاک، نویسندهٔ دیگر قرن نوزدهم است؛ به‌طوری‌که دو رمان بسیار مشهور وی، مادام بواری و تربیت احساسات، به ترتیب از زن سی ساله و زنبق درهٔ بالزاک الهام می‌گیرند.	2025-08-19 13:26:04.63661	2025-08-19 13:26:04.63661	1
19	ویکتور هوگو	فرانسه	ویکتور-ماری هوگو (فرانسوی: Victor-Marie Hugo‎؛ ۲۶ فوریه ۱۸۰۲ – ۲۲ مه ۱۸۸۵) نویسنده و سیاستمدار فرانسوی پیرو رمانتیسم بود. او طی دوران کاری ادبی‌اش که بیش از شش دهه به طول انجامید، در ژانرها و اشکال مختلف نوشت.\r\n\r\nآثار او به بسیاری از اندیشه‌های سیاسی و هنری رایج در زمان خویش اشاره کرده و بازگویندهٔ تاریخ معاصر فرانسه است. از برجسته‌ترین آثار او بینوایان، گوژپشت نوتردام، کارگران دریا، مردی که می‌خندد و آخرین روز یک محکوم است. مشهورترین کارهای هوگو در خارج از فرانسه بینوایان و گوژپشت نوتردام است و در فرانسه وی را با مجموعه اشعارش هم می‌شناسند.	2025-08-19 13:28:37.615597	2025-08-19 13:28:37.615597	1
20	جی. دی. سلینجر	آمریکا	جروم دیوید سَلینجر (به آلمانی: Jerome David Salinger) (زادهٔ ۱ ژانویهٔ ۱۹۱۹ – درگذشتهٔ ۲۷ ژانویهٔ ۲۰۱۰) نویسندهٔ معاصر آمریکایی بود. رمان‌های پرطرفدار وی، مانند ناطور دشت در نقد جامعهٔ مدرن غرب و خصوصاً آمریکا نوشته شده‌اند. سلینجر بیشتر با حروف ابتداییِ نام خود «جی. دی. سلینجر» معروف است.	2025-08-19 13:32:14.020184	2025-08-19 13:32:14.020184	1
\.


--
-- Data for Name: books; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.books (id, title, isbn, author_id, total_copies, available_copies, created_at, updated_at, version, publisher, language, summary, translator, cover_image_url) FROM stdin;
134	صد سال تنهایی	9780060883287	14	10	10	2025-08-19 09:29:11.535657	2025-08-19 09:29:11.535657	1	Editorial Sudamericana	اسپانیایی	صد سال تنهایی (به اسپانیایی: Cien años de soledad) نام رمانی به زبان اسپانیایی نوشته گابریل گارسیا مارکز که چاپ نخست آن در سال ۱۹۶۷ در آرژانتین با تیراژ ۸۰۰۰ نسخه منتشر شد. تمام نسخه‌های چاپ اول صد سال تنهایی به زبان اسپانیایی در همان هفته نخست کاملاً به فروش رفت. در ۳۰ سالی که از نخستین چاپ این کتاب گذشت بیش از ۳۰ میلیون نسخه از آن در سراسر جهان به فروش رفته و به بیش از ۳۰ زبان ترجمه شده‌است. جایزه نوبل ادبیات ۱۹۸۲ به گابریل گارسیا مارکز به خاطر خلق این اثر تعلق گرفت.[۱]\r\n\r\n		/static/covers/Cien_años_de_soledad_(book_cover,_1967).jpg
131	1984	9780451524935	11	11	11	2025-08-19 09:08:51.02882	2025-08-19 09:13:19.45668	2	Secker & Warburg	انگلیسی	۱۹۸۴ نام کتاب مشهوری از جورج اورول که در سال ۱۹۴۹ منتشر شده‌است.[۱] این کتاب بیانیهٔ سیاسی شاخصی در رده‌ی نظام‌های تمامیت‌خواه (توتالیتر) شمرده می‌شود. ۱۹۸۴ کتابی پادآرمانشهری به‌شمار می‌آید.		/static/covers/1984.jpg
132	کشتن مرغ مقلد	9780061120084	12	5	5	2025-08-19 09:21:00.103462	2025-08-19 09:21:00.103462	1	J.B. Lippincott & Co.	انگلیسی	کشتن مرغ مقلد یا کشتن مرغ مینا (به انگلیسی: To Kill a Mockingbird) رمانی نوشته هارپر لی، نویسنده آمریکایی در سال ۱۹۶۰ میلادی است.[۱]\r\n\r\nوی برای نوشتن این رمان در سال ۱۹۶۴ میلادی، جایزه پولیتزر را به دست آورد. از زمان اولین انتشار تاکنون، بیش از ۴۰ میلیون نسخه از این کتاب به فروش رفته و به بیش از ۴۰ زبان بین‌المللی ترجمه شده‌است.[۲] تیراژ نخستین چاپ رُمان کشتن مرغ مقلد ۵ هزار نسخه بود.[۲] با وجود برخورد با مسائل جدی مانند تجاوز به عنف و نابرابری نژادی، این رمان برای گرما و طنز آن مشهور است. پدر راوی، اتیکاس فینچ، به عنوان یک قهرمان اخلاقی برای بسیاری از خوانندگان و به عنوان یک مدل از یکپارچگی برای وکلا نشان داده شده‌است. به عنوان یک رمان گوتیک جنوبی و رمان تربیتی، تم اولیه کشتن مرغ مقلد شامل بی عدالتی نژادی و کشتار بی گناهی است. محققان اشاره کرده‌اند که لی همچنین مسائل مربوط به طبقات اجتماعی، شجاعت، محبت و نقش‌های جنسیتی در جنوب آمریکا را به خوبی به مخاطب منتقل کرده‌است.[۳]		/static/covers/To_Kill_a_Mockingbird_(first_edition_cover).jpg
135	جنایت و مکافات	9780140449136	15	8	8	2025-08-19 09:32:51.15238	2025-08-19 09:32:51.15238	1	The Russian Messenger	روسی	این کتاب داستان دانشجویی به نام راسکولْنیکُف را روایت می‌کند که به‌خاطر اصول مرتکب قتل می‌شود. بنابر انگیزه‌های پیچیده‌ای که حتی خود او از تحلیلشان عاجز است، زن رباخوار یهودی را همراه با خواهرش که غیرمنتظره به هنگام وقوع قتل در صحنه حاضر می‌شود، می‌کُشد و پس از قتل خود را ناتوان از خرج کردن پول و جواهراتی که برداشته می‌بیند و آن‌ها را پنهان می‌کند. بعد از چند روز بیماری و بستری شدن در خانه راسکولنیکف هرکس را که می‌بیند می‌پندارد به او مظنون است و با این افکار کارش به جنون می‌رسد. در این بین او عاشق سونیا، دختری که به‌خاطر مشکلات مالی خانواده‌اش دست به تن‌فروشی زده بود، می‌شود. داستایفسکی این رابطه را به نشانهٔ مِهر خداوندی به انسان خطاکار استفاده کرده‌است و همان عشق، نیروی رستگاری‌بخش است. البته راسکولنیکف بعد از اقرار به گناه و زندانی شدن در سیبری به این حقیقت رسید.\r\n\r\n		/static/covers/crimeandpunishment.png
133	گتسبی بزرگ	9780743273565	13	5	4	2025-08-19 09:25:30.796564	2025-08-19 10:33:42.209989	2	Charles Scribner's Sons	انگلیسی	گتسبی بزرگ (انگلیسی: The Great Gatsby) رمانی نوشتهٔ نویسندهٔ آمریکایی، اف. اسکات فیتزجرالد است که در سال ۱۹۲۵ منتشر شد. این رمان که داستان آن در عصر جاز در لانگ آیلند، در نزدیکی نیویورک رخ می‌دهد، روابط متقابل نیک کاراوی، راوی اول‌شخص داستان با جی گتسبی، میلیونر مرموز و فکر همیشگی گتسبی برای پیوستن دوباره به معشوقهٔ سابقش، دیزی بیوکنن را شرح می‌دهد.		/static/covers/The_Great_Gatsby_Cover_1925_Retouched.jpg
138	برادران کارامازوف	9780374528379	15	4	3	2025-08-19 09:54:50.463692	2025-08-19 10:22:25.776767	3	The Russian Messenger	روسی	برادران کارامازوف (روسی: Братья Карамазовы، تلفظ [ˈbratʲjə kərɐˈmazəvɨ]) رمانی از فیودور داستایفسکی — نویسندهٔ روس — است که نخستین بار در سال‌های ۸۰–۱۸۷۹ در نشریهٔ پیام‌آور روسی به‌صورت پاورقی منتشر شد. نوشته‌های نویسنده نشان می‌دهند که این کتاب قرار بود قسمت اول از مجموعه‌ای بزرگتر با نام زندگی یک گناهکار بزرگ باشد؛ ولی قسمت دوم کتاب هیچگاه نوشته نشد، چون داستایفسکی چهار ماه بعد از پایان انتشار کتاب درگذشت. این کتاب در چهار جلد منتشر شد و مثل بیشتر آثار داستایفسکی، یک رمان جنایی است.		/static/covers/karamazov.jpg
136	جنگ و صلح	9780199232765	16	5	4	2025-08-19 09:36:02.578096	2025-08-19 10:34:38.122712	4	The Russian Messenger	روسی	جنگ و صلح (به روسی: Война́ и мир یا Voyná i mir) نام رمان مشهور لئو تولستوی، نویسنده بزرگ روس است.\r\n\r\nتولستوی کتاب جنگ و صلح را در سال ۱۸۶۹ میلادی نوشت. این کتاب یکی از بزرگ‌ترین آثار ادبیات روسی و از مهم‌ترین رمان‌های ادبیات جهان به‌شمار می‌رود.[۱] در این رمان طولانی بیش از ۵۸۰ شخصیت با دقت توصیف شده‌اند و یکی از معتبرترین منابع تحقیق و بررسی در تاریخ سیاسی و اجتماعی سده نوزدهم امپراتوری روسیه است که به شرح مقاومت روس‌ها در برابر حملهٔ ارتش فرانسه به رهبری ناپلئون بناپارت می‌پردازد. منتقدان ادبی آن را یکی از بزرگ‌ترین رمان‌های جهان می‌دانند.[۲] این رمان، زندگی اجتماعی و سرگذشت چهار خانواده اشرافی[۳] روس را در دوران جنگ‌های روسیه و فرانسه در سال‌های ۱۸۰۵ تا ۱۸۱۴ به تصویر می‌کشد.[۴]\r\n\r\n		/static/covers/war-and-peace.jpg
139	مادام بوواری	9780140449129	18	5	5	2025-08-19 09:56:58.425995	2025-08-19 11:59:41.217422	3	Revue de Paris	فرانسوی	مادام بوواری (به فرانسوی: Madame Bovary) نخستین اثر گوستاو فلوبر، نویسنده نامدار فرانسوی است که یکی از برجسته‌ترین آثار او به‌شمار می‌آید.\r\n\r\n		/static/covers/Madame_Bovary_1857_(hi-res).jpg
140	بینوایان	9780451419439	19	10	10	2025-08-19 10:00:37.87746	2025-08-19 10:00:37.87746	1	A. Lacroix, Verboeckhoven & Cie	فرانسوی	بینوایان (به فرانسوی: Les Misérables) یک رمان تاریخی فرانسوی نوشتهٔ ویکتور هوگو است که اولین بار در سال ۱۸۶۲ منتشر شد و به‌عنوان یکی از بزرگ‌ترین رمانهای قرن ۱۹م شناخته می‌شود. رمان از آغاز شورش ژوئن در پاریسِ سال ۱۸۱۵ تا به‌ثمر رسیدن آن در سال ۱۸۳۲، زندگی شخصیت‌های مختلف، به‌ویژه زندانی آزادشده‌ای به نام ژان والژان را روایت می‌کند.\r\n\r\nاین رمان با بررسی ماهیت قانون و بخشش، تاریخ فرانسه، معماری و طراحی شهریِ پاریس، سیاست‌ها، فلسفه اخلاق، ضد اخلاقیات، قضاوت‌ها، مذهب، نوع و ماهیت عشق را شرح می‌دهد. بینوایان از طریق فیلم، برنامه‌های تلویزیونی و تئاتر به اقبال عمومی فراوان دست یافت. فیلم‌هایی مثل بینوایان (موزیکال)، بینوایان (فیلم ۱۹۹۸)، بینوایان (فیلم ۲۰۱۲) و بینوایان (مینی‌سریال ۲۰۱۸) از جملهٔ آثار هنری اقتباسی از این رمان هستند.		/static/covers/Monsieur_Madeleine_par_Gustave_Brion.jpg
141	ناطور دشت	9780316769488	20	10	9	2025-08-19 10:03:32.446068	2025-08-19 10:16:41.891151	2	Little, Brown and Company	انگلیسی	ناطورِ دشت یا ناتورِ دشت (انگلیسی: The Catcher in the Rye) نام رمانی اثر نویسندهٔ آمریکایی جروم دیوید سلینجر است که در ابتدا به صورت دنباله‌دار در سال‌های ۱۹۴۵–۴۶ انتشار یافت.[۱] این رمان کلاسیک که در اصل برای مخاطب بزرگسال منتشر شده بود، به دلیل درون‌مایه طغیان‌گری و عصبانیت نوجوان داستان، مورد توجه بسیاری از نوجوانان قرار گرفت.		/static/covers/Catcher-in-the-rye.jpg
137	دن کیشوت	9780060934347	17	10	9	2025-08-19 09:39:16.872562	2025-08-19 10:21:06.512534	2	Francisco de Robles	اسپانیایی	دُن کیشوت زندگی فردی را روایت می‌کند که دچار توهم است و وقت خود را با خواندن آثار ممنوعه می‌گذراند. در زمان روایت داستان، نوشتن و خواندن آثاری که به شوالیه‌ها می‌پرداخت، ممنوع بود و شخصیت اصلی داستان، خود را در نقش یکی از همین شوالیه‌ها تصور می‌کند و دشمنانی فرضی در برابر خود می‌بیند که اغلب کوه‌ها و درخت‌ها هستند. دُن کیشوت پهلوانی خیالی و بی‌دست‌وپاست که خود را شکست‌ناپذیر می‌پندارد.\r\n\r\nاو با همراهی خدمتکارش، سانچو پانزا، به سفرهایی طولانی می‌رود و در میانهٔ همین سفرهاست که کارهایی عجیب و غریب از وی سر می‌زند. او که هدفی جز نجات مردم از ظلم و استبداد حاکمان ظالم ندارد، نگاهی تخیلی به اطرافش دارد و همه چیز را در قالب ابزار جنگی می‌بیند.[۱]\r\n\r\n		/static/covers/El_ingenioso_hidalgo_don_Quijote_de_la_Mancha.jpg
\.


--
-- Data for Name: categories; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.categories (id, slug, name, created_at, updated_at, version) FROM stdin;
\.


--
-- Data for Name: category_books; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.category_books (book_id, category_id) FROM stdin;
\.


--
-- Data for Name: loans; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.loans (id, book_id, member_id, borrow_date, due_date, return_date, status, created_at, updated_at, version) FROM stdin;
60	141	54	2025-08-19 00:00:00	2025-09-06 00:00:00	\N	borrowed	2025-08-19 13:46:41.891252	2025-08-19 13:46:41.891252	1
61	137	59	2025-08-19 00:00:00	2025-08-22 00:00:00	\N	borrowed	2025-08-19 13:51:06.512632	2025-08-19 13:51:06.512632	1
62	138	62	2025-08-19 00:00:00	2025-08-23 00:00:00	\N	borrowed	2025-08-19 13:52:25.776842	2025-08-19 13:52:25.776842	1
63	133	61	2025-08-19 00:00:00	2025-08-30 00:00:00	\N	borrowed	2025-08-19 14:03:42.21008	2025-08-19 14:03:42.21008	1
64	136	61	2025-08-19 00:00:00	2025-08-23 00:00:00	\N	borrowed	2025-08-19 14:04:38.122854	2025-08-19 14:04:38.122854	1
\.


--
-- Data for Name: members; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.members (id, full_name, email, phone_number, joined_at, status, created_at, updated_at, version, national_id) FROM stdin;
54	علی رضایی	ali.rezaei@example.com	09124567890	2025-08-19 10:07:11.772007	active	2025-08-19 13:37:11.772398	2025-08-19 10:07:26.18013	2	00546632
55	زهرا محمدی	zahra.mohammadi@example.com	09357442121	2025-08-19 10:08:00.468473	active	2025-08-19 13:38:00.468597	2025-08-19 13:38:00.468597	1	001234323
56	مهدی کریمی	mehdi.karimi@example.com	09132238541	2025-08-19 10:08:49.665684	active	2025-08-19 13:38:49.665868	2025-08-19 13:38:49.665868	1	009327711
57	سارا احمدی	sarah.ahmadi@example.com	09393228981	2025-08-19 10:09:18.473166	active	2025-08-19 13:39:18.473293	2025-08-19 13:39:18.473293	1	001923342
58	رضا مرادی	reza.moradi@example.com	09199283310	2025-08-19 10:10:35.30732	active	2025-08-19 13:40:35.30741	2025-08-19 13:40:35.30741	1	0032831234
59	فاطمه حسینی	fateme.hosseini@example.com	09129343219	2025-08-19 10:11:39.999046	active	2025-08-19 13:41:39.999152	2025-08-19 13:41:39.999152	1	۰۰۹۴۴۸۳۲۱
60	محمدرضا کاظمی	mreza.kazemi@example.com	09338921293	2025-08-19 10:12:36.830702	active	2025-08-19 13:42:36.830795	2025-08-19 13:42:36.830795	1	003873581
61	نگار صفایی	negar.saffaei@example.com	09358329111	2025-08-19 10:13:13.243648	active	2025-08-19 13:43:13.243733	2025-08-19 13:43:13.243733	1	009542376
62	امیر نادری	amir.naderi@example.com	09198432241	2025-08-19 10:13:45.530041	active	2025-08-19 13:43:45.530153	2025-08-19 13:43:45.530153	1	003449322
\.


--
-- Data for Name: staff; Type: TABLE DATA; Schema: public; Owner: -
--

COPY public.staff (id, full_name, phone_number, email, role, password_hash, created_at, updated_at, last_login, version, status) FROM stdin;
10	محمدمهدی بوالحسنی	09350000000	admin@lms.org	admin	$2a$10$f7RkrRwbgcPxlZDfzlMcTu3duoTDkZj5QPxHxz4hV/flsZ9eGQEu6	2025-08-10 02:26:52.545505	2025-08-20 16:24:19.921053	2025-08-19 21:00:02.053924	18	active
\.


--
-- Name: activity_log_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.activity_log_id_seq', 77, true);


--
-- Name: authors_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.authors_id_seq', 20, true);


--
-- Name: books_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.books_id_seq', 141, true);


--
-- Name: category_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.category_id_seq', 9, true);


--
-- Name: loans_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.loans_id_seq', 65, true);


--
-- Name: members_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.members_id_seq', 62, true);


--
-- Name: staff_id_seq; Type: SEQUENCE SET; Schema: public; Owner: -
--

SELECT pg_catalog.setval('public.staff_id_seq', 12, true);


--
-- Name: activity_logs activity_log_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.activity_logs
    ADD CONSTRAINT activity_log_pkey PRIMARY KEY (id);


--
-- Name: authors authors_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.authors
    ADD CONSTRAINT authors_pkey PRIMARY KEY (id);


--
-- Name: books books_isbn_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.books
    ADD CONSTRAINT books_isbn_key UNIQUE (isbn);


--
-- Name: books books_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.books
    ADD CONSTRAINT books_pkey PRIMARY KEY (id);


--
-- Name: category_books category_books_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.category_books
    ADD CONSTRAINT category_books_pkey PRIMARY KEY (book_id, category_id);


--
-- Name: categories category_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT category_pkey PRIMARY KEY (id);


--
-- Name: categories category_slug_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.categories
    ADD CONSTRAINT category_slug_key UNIQUE (slug);


--
-- Name: loans loans_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.loans
    ADD CONSTRAINT loans_pkey PRIMARY KEY (id);


--
-- Name: members members_email_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.members
    ADD CONSTRAINT members_email_key UNIQUE (email);


--
-- Name: members members_national_id_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.members
    ADD CONSTRAINT members_national_id_key UNIQUE (national_id);


--
-- Name: members members_phone_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.members
    ADD CONSTRAINT members_phone_key UNIQUE (phone_number);


--
-- Name: members members_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.members
    ADD CONSTRAINT members_pkey PRIMARY KEY (id);


--
-- Name: staff staff_email_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.staff
    ADD CONSTRAINT staff_email_key UNIQUE (email);


--
-- Name: staff staff_phone_number_key; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.staff
    ADD CONSTRAINT staff_phone_number_key UNIQUE (phone_number);


--
-- Name: staff staff_pkey; Type: CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.staff
    ADD CONSTRAINT staff_pkey PRIMARY KEY (id);


--
-- Name: unique_active_loan; Type: INDEX; Schema: public; Owner: -
--

CREATE UNIQUE INDEX unique_active_loan ON public.loans USING btree (member_id, book_id) WHERE (status = 'borrowed'::text);


--
-- Name: active_members_view _RETURN; Type: RULE; Schema: public; Owner: -
--

CREATE OR REPLACE VIEW public.active_members_view AS
 SELECT members.id,
    members.full_name,
    count(loans.id) AS loan_count
   FROM (public.members
     JOIN public.loans ON ((members.id = loans.member_id)))
  GROUP BY members.id
  ORDER BY (count(loans.id)) DESC;


--
-- Name: popular_books_view _RETURN; Type: RULE; Schema: public; Owner: -
--

CREATE OR REPLACE VIEW public.popular_books_view AS
 SELECT books.id,
    books.title,
    count(loans.id) AS loan_count
   FROM (public.books
     JOIN public.loans ON ((books.id = loans.book_id)))
  GROUP BY books.id
  ORDER BY (count(loans.id)) DESC;


--
-- Name: low_stock_high_demand_books_view _RETURN; Type: RULE; Schema: public; Owner: -
--

CREATE OR REPLACE VIEW public.low_stock_high_demand_books_view AS
 SELECT b.id,
    b.title,
    b.total_copies,
    b.available_copies,
    count(l.id) AS loan_count
   FROM (public.books b
     JOIN public.loans l ON ((l.book_id = b.id)))
  WHERE ((l.borrow_date >= (CURRENT_DATE - '6 mons'::interval)) AND (b.available_copies <= 2))
  GROUP BY b.id, b.title
  ORDER BY (count(l.id)) DESC;


--
-- Name: popular_authors_view _RETURN; Type: RULE; Schema: public; Owner: -
--

CREATE OR REPLACE VIEW public.popular_authors_view AS
 SELECT a.id,
    a.full_name,
    count(l.id) AS total_loans
   FROM ((public.authors a
     JOIN public.books b ON ((b.author_id = a.id)))
     JOIN public.loans l ON ((l.book_id = b.id)))
  GROUP BY a.id
  ORDER BY (count(l.id)) DESC;


--
-- Name: loans trg_prevent_returned_loan_insertion; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER trg_prevent_returned_loan_insertion BEFORE INSERT ON public.loans FOR EACH ROW WHEN ((new.status = 'returned'::text)) EXECUTE FUNCTION public.prevent_insert_loan_as_returend();


--
-- Name: books trg_set_available_copies_to_total; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER trg_set_available_copies_to_total BEFORE INSERT ON public.books FOR EACH ROW EXECUTE FUNCTION public.set_available_copies_to_total();


--
-- Name: loans trg_sync_available_copies; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER trg_sync_available_copies AFTER INSERT OR DELETE OR UPDATE ON public.loans FOR EACH ROW EXECUTE FUNCTION public.update_available_copies();


--
-- Name: books trg_sync_available_copies_on_totalcopy_update; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER trg_sync_available_copies_on_totalcopy_update BEFORE UPDATE OF total_copies ON public.books FOR EACH ROW WHEN ((old.total_copies <> new.total_copies)) EXECUTE FUNCTION public.update_available_copies_on_totalcopy_update();


--
-- Name: authors trg_update_update_fields; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER trg_update_update_fields BEFORE UPDATE ON public.authors FOR EACH ROW EXECUTE FUNCTION public.update_update_fields();


--
-- Name: books trg_update_update_fields; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER trg_update_update_fields BEFORE UPDATE ON public.books FOR EACH ROW EXECUTE FUNCTION public.update_update_fields();


--
-- Name: categories trg_update_update_fields; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER trg_update_update_fields BEFORE UPDATE ON public.categories FOR EACH ROW EXECUTE FUNCTION public.update_update_fields();


--
-- Name: loans trg_update_update_fields; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER trg_update_update_fields BEFORE UPDATE ON public.loans FOR EACH ROW EXECUTE FUNCTION public.update_update_fields();


--
-- Name: members trg_update_update_fields; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER trg_update_update_fields BEFORE UPDATE ON public.members FOR EACH ROW EXECUTE FUNCTION public.update_update_fields();


--
-- Name: staff trg_update_update_fields; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER trg_update_update_fields BEFORE UPDATE ON public.staff FOR EACH ROW EXECUTE FUNCTION public.update_update_fields();


--
-- Name: books trg_validate_available_copies_on_update; Type: TRIGGER; Schema: public; Owner: -
--

CREATE TRIGGER trg_validate_available_copies_on_update BEFORE UPDATE OF available_copies ON public.books FOR EACH ROW WHEN ((new.available_copies <> old.available_copies)) EXECUTE FUNCTION public.validate_available_copies_update();


--
-- Name: books books_author_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.books
    ADD CONSTRAINT books_author_fkey FOREIGN KEY (author_id) REFERENCES public.authors(id) ON DELETE SET NULL;


--
-- Name: category_books category_books_book_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.category_books
    ADD CONSTRAINT category_books_book_id_fkey FOREIGN KEY (book_id) REFERENCES public.books(id) ON DELETE CASCADE;


--
-- Name: category_books category_books_category_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.category_books
    ADD CONSTRAINT category_books_category_id_fkey FOREIGN KEY (category_id) REFERENCES public.categories(id) ON DELETE CASCADE;


--
-- Name: loans loans_book_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.loans
    ADD CONSTRAINT loans_book_id_fkey FOREIGN KEY (book_id) REFERENCES public.books(id) ON DELETE RESTRICT;


--
-- Name: loans loans_member_id_fkey; Type: FK CONSTRAINT; Schema: public; Owner: -
--

ALTER TABLE ONLY public.loans
    ADD CONSTRAINT loans_member_id_fkey FOREIGN KEY (member_id) REFERENCES public.members(id) ON DELETE CASCADE;


--
-- PostgreSQL database dump complete
--

