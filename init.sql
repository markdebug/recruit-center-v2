-- 清理现有数据（按依赖关系逆序删除）
DELETE FROM t_rc_job_apply;
DELETE FROM t_rc_job_statistics;
DELETE FROM t_rc_job_favorite;
DELETE FROM t_rc_resume_interaction;
DELETE FROM t_rc_resume_attachment;
DELETE FROM t_rc_resume_project;
DELETE FROM t_rc_resume_work_experience;
DELETE FROM t_rc_resume_education;
DELETE FROM t_rc_resume;
DELETE FROM t_rc_job;
DELETE FROM t_rc_notification;
DELETE FROM t_rc_notification_template;
DELETE FROM t_rc_dict;

-- 重置序列
ALTER SEQUENCE t_rc_resume_id_seq RESTART WITH 1;
ALTER SEQUENCE t_rc_resume_education_id_seq RESTART WITH 1;
ALTER SEQUENCE t_rc_resume_work_experience_id_seq RESTART WITH 1;
ALTER SEQUENCE t_rc_resume_project_id_seq RESTART WITH 1;
ALTER SEQUENCE t_rc_resume_attachment_id_seq RESTART WITH 1;
ALTER SEQUENCE t_rc_resume_interaction_id_seq RESTART WITH 1;
ALTER SEQUENCE t_rc_job_id_seq RESTART WITH 1;
ALTER SEQUENCE t_rc_job_apply_id_seq RESTART WITH 1;
ALTER SEQUENCE t_rc_job_statistics_id_seq RESTART WITH 1;
ALTER SEQUENCE t_rc_job_favorite_id_seq RESTART WITH 1;
ALTER SEQUENCE t_rc_notification_id_seq RESTART WITH 1;
ALTER SEQUENCE t_rc_notification_template_id_seq RESTART WITH 1;
ALTER SEQUENCE t_rc_dict_id_seq RESTART WITH 1;

-- 插入字典数据 (修正列名)
INSERT INTO t_rc_dict (id, parent_id, category, code, name, value, sort, status, remarks, created_at, updated_at) VALUES
(1, 0, 'job_category', 'tech', '技术岗位', 'tech', 1, 1, '技术相关职位', NOW(), NOW()),
(2, 0, 'job_category', 'product', '产品岗位', 'product', 2, 1, '产品相关职位', NOW(), NOW()),
(3, 0, 'job_category', 'design', '设计岗位', 'design', 3, 1, '设计相关职位', NOW(), NOW()),
(4, 0, 'job_category', 'marketing', '市场岗位', 'marketing', 4, 1, '市场营销职位', NOW(), NOW()),
(5, 0, 'job_category', 'sales', '销售岗位', 'sales', 5, 1, '销售相关职位', NOW(), NOW()),
(6, 0, 'education_level', 'bachelor', '本科', 'bachelor', 1, 1, '本科学历', NOW(), NOW()),
(7, 0, 'education_level', 'master', '硕士', 'master', 2, 1, '硕士学历', NOW(), NOW()),
(8, 0, 'education_level', 'phd', '博士', 'phd', 3, 1, '博士学历', NOW(), NOW()),
(9, 0, 'experience_level', '1-3', '1-3年', '1-3', 1, 1, '1到3年经验', NOW(), NOW()),
(10, 0, 'experience_level', '3-5', '3-5年', '3-5', 2, 1, '3到5年经验', NOW(), NOW());

-- 插入通知模板数据 (修正列名)
INSERT INTO t_rc_notification_template (id, code, title, content, type, user_types, channels, is_active, create_time, update_time) VALUES
(1, 'job_apply_success', '职位投递成功', '您已成功投递{{jobName}}职位，请耐心等待HR回复。', 1, '[1]', 1, true, NOW(), NOW()),
(2, 'job_apply_review', '简历已通过筛选', '恭喜！您投递的{{jobName}}职位简历已通过初筛，HR将尽快与您联系。', 2, '[1]', 1, true, NOW(), NOW()),
(3, 'job_apply_reject', '投递结果通知', '很遗憾，您投递的{{jobName}}职位暂未通过筛选，感谢您的关注。', 2, '[1]', 1, true, NOW(), NOW()),
(4, 'interview_invite', '面试邀请', '您好！{{companyName}}邀请您参加{{jobName}}职位的面试，时间：{{interviewTime}}。', 3, '[1]', 3, true, NOW(), NOW()),
(5, 'offer_send', 'Offer通知', '恭喜！{{companyName}}向您发送了{{jobName}}职位的Offer，请及时查看。', 2, '[1]', 1, true, NOW(), NOW()),
(6, 'resume_view', '简历被查看', '您的简历被{{companyName}}查看了，快去看看吧！', 4, '[1]', 1, true, NOW(), NOW()),
(7, 'job_expire_remind', '职位即将到期', '您发布的{{jobName}}职位将在3天后到期，请及时续期。', 4, '[2]', 1, true, NOW(), NOW()),
(8, 'system_maintain', '系统维护通知', '系统将于{{maintainTime}}进行维护，预计耗时{{duration}}。', 4, '[1,2,3,4]', 1, true, NOW(), NOW()),
(9, 'new_job_recommend', '为您推荐新职位', '根据您的简历，为您推荐{{jobName}}职位，快来看看吧！', 1, '[1]', 1, true, NOW(), NOW()),
(10, 'profile_complete', '完善简历信息', '完善您的简历信息，获得更多面试机会！', 4, '[1]', 1, true, NOW(), NOW());

-- 插入简历数据
INSERT INTO t_rc_resume (id, user_id, name, avatar, gender, birthday, phone, email, location, experience, job_status, expected_job, expected_city, expected_salary, introduction, skills, share_token, access_status, working_status, status, created_at, updated_at) VALUES
(1, 1001, '张三', 'https://example.com/avatar1.jpg', 1, '1990-05-15', 'encrypted_phone_1', 'encrypted_email_1', 'encrypted_location_1', 5, 1, 'Go开发工程师', '北京', '20K-30K', '5年Go开发经验，熟悉微服务架构', 'Go,Docker,Kubernetes,MySQL', 'share_token_1', 2, 1, 1, NOW(), NOW()),
(2, 1002, '李四', 'https://example.com/avatar2.jpg', 2, '1992-08-20', 'encrypted_phone_2', 'encrypted_email_2', 'encrypted_location_2', 3, 1, 'Java开发工程师', '上海', '15K-25K', '3年Java开发经验，熟悉Spring框架', 'Java,Spring,Redis,MongoDB', 'share_token_2', 2, 1, 1, NOW(), NOW()),
(3, 1003, '王五', 'https://example.com/avatar3.jpg', 1, '1988-12-10', 'encrypted_phone_3', 'encrypted_email_3', 'encrypted_location_3', 7, 2, '技术经理', '深圳', '30K-40K', '7年开发经验，3年管理经验', 'Python,Django,PostgreSQL,团队管理', 'share_token_3', 2, 1, 1, NOW(), NOW()),
(4, 1004, '赵六', 'https://example.com/avatar4.jpg', 2, '1995-03-25', 'encrypted_phone_4', 'encrypted_email_4', 'encrypted_location_4', 2, 1, '前端开发工程师', '杭州', '12K-18K', '2年前端开发经验，熟悉Vue.js', 'Vue.js,React,JavaScript,CSS', 'share_token_4', 2, 1, 1, NOW(), NOW()),
(5, 1005, '钱七', 'https://example.com/avatar5.jpg', 1, '1987-11-08', 'encrypted_phone_5', 'encrypted_email_5', 'encrypted_location_5', 8, 1, '架构师', '广州', '35K-50K', '8年开发经验，负责大型系统架构设计', 'Java,微服务,分布式,架构设计', 'share_token_5', 2, 1, 1, NOW(), NOW()),
(6, 1006, '孙八', 'https://example.com/avatar6.jpg', 2, '1993-07-14', 'encrypted_phone_6', 'encrypted_email_6', 'encrypted_location_6', 4, 1, '产品经理', '成都', '18K-28K', '4年产品经验，擅长用户体验设计', '产品设计,用户研究,原型设计,需求分析', 'share_token_6', 2, 1, 1, NOW(), NOW()),
(7, 1007, '周九', 'https://example.com/avatar7.jpg', 1, '1991-04-30', 'encrypted_phone_7', 'encrypted_email_7', 'encrypted_location_7', 6, 1, 'DevOps工程师', '西安', '22K-32K', '6年运维经验，熟悉云原生技术', 'Docker,Kubernetes,AWS,Jenkins', 'share_token_7', 2, 1, 1, NOW(), NOW()),
(8, 1008, '吴十', 'https://example.com/avatar8.jpg', 2, '1994-09-18', 'encrypted_phone_8', 'encrypted_email_8', 'encrypted_location_8', 3, 1, 'UI设计师', '南京', '15K-22K', '3年UI设计经验，擅长移动端设计', 'Sketch,Figma,Adobe Creative Suite,UI设计', 'share_token_8', 2, 1, 1, NOW(), NOW()),
(9, 1009, '郑十一', 'https://example.com/avatar9.jpg', 1, '1989-01-22', 'encrypted_phone_9', 'encrypted_email_9', 'encrypted_location_9', 9, 1, '数据科学家', '武汉', '25K-35K', '9年数据分析经验，擅长机器学习', 'Python,R,TensorFlow,机器学习,数据挖掘', 'share_token_9', 2, 1, 1, NOW(), NOW()),
(10, 1010, '冯十二', 'https://example.com/avatar10.jpg', 2, '1996-06-05', 'encrypted_phone_10', 'encrypted_email_10', 'encrypted_location_10', 1, 1, '测试工程师', '青岛', '10K-15K', '1年测试经验，熟悉自动化测试', 'Selenium,JMeter,Python,自动化测试', 'share_token_10', 2, 1, 1, NOW(), NOW());

-- 插入教育经历数据
INSERT INTO t_rc_resume_education (id, resume_id, school, major, degree, start_time, end_time, created_at, updated_at) VALUES
(1, 1, '清华大学', '计算机科学与技术', '本科', '2009-09-01', '2013-06-30', NOW(), NOW()),
(2, 1, '清华大学', '软件工程', '硕士', '2013-09-01', '2016-06-30', NOW(), NOW()),
(3, 2, '北京大学', '计算机科学与技术', '本科', '2011-09-01', '2015-06-30', NOW(), NOW()),
(4, 3, '上海交通大学', '软件工程', '本科', '2007-09-01', '2011-06-30', NOW(), NOW()),
(5, 4, '浙江大学', '计算机科学与技术', '本科', '2014-09-01', '2018-06-30', NOW(), NOW()),
(6, 5, '华中科技大学', '计算机科学与技术', '本科', '2005-09-01', '2009-06-30', NOW(), NOW()),
(7, 6, '西安电子科技大学', '电子信息工程', '本科', '2012-09-01', '2016-06-30', NOW(), NOW()),
(8, 7, '北京理工大学', '计算机科学与技术', '本科', '2010-09-01', '2014-06-30', NOW(), NOW()),
(9, 8, '中央美术学院', '视觉传达设计', '本科', '2013-09-01', '2017-06-30', NOW(), NOW()),
(10, 9, '复旦大学', '数学与应用数学', '本科', '2008-09-01', '2012-06-30', NOW(), NOW());

-- 插入工作经历数据
INSERT INTO t_rc_resume_work_experience (id, resume_id, company_name, position, department, start_time, end_time, description, achievement, created_at, updated_at) VALUES
(1, 1, '阿里巴巴', 'Go开发工程师', '技术部', '2016-07-01', '2021-12-31', '负责电商后端服务开发', '优化订单处理性能提升50%', NOW(), NOW()),
(2, 1, '字节跳动', '高级Go开发工程师', '基础架构部', '2022-01-01', '2024-01-01', '负责微服务架构设计与实现', '设计的微服务架构支撑日活千万用户', NOW(), NOW()),
(3, 2, '腾讯', 'Java开发工程师', '游戏部', '2018-07-01', '2021-08-31', '负责游戏后端开发', '开发的系统支持百万并发用户', NOW(), NOW()),
(4, 3, '百度', '高级Java开发工程师', '搜索部', '2015-07-01', '2019-12-31', '负责搜索引擎后端开发', '参与核心搜索算法优化', NOW(), NOW()),
(5, 3, '美团', '技术经理', '到店事业群', '2020-01-01', '2024-01-01', '负责团队管理和技术架构', '带领20人团队完成多个重要项目', NOW(), NOW()),
(6, 4, '网易', '前端开发工程师', '游戏部', '2020-07-01', '2022-12-31', '负责游戏前端界面开发', '开发的H5游戏获得千万用户', NOW(), NOW()),
(7, 5, '华为', '软件工程师', '云计算部', '2011-07-01', '2016-06-30', '负责云平台开发', '参与华为云核心组件开发', NOW(), NOW()),
(8, 6, '小米', '产品经理', '生态链部', '2018-07-01', '2022-12-31', '负责智能硬件产品规划', '主导3款爆品产品从0到1', NOW(), NOW()),
(9, 7, '京东', 'DevOps工程师', '技术部', '2016-07-01', '2022-12-31', '负责持续集成和部署', '构建的CI/CD系统提升发布效率70%', NOW(), NOW()),
(10, 8, '滴滴', 'UI设计师', '设计部', '2019-07-01', '2022-12-31', '负责移动端UI设计', '设计的界面获得用户体验大奖', NOW(), NOW());

-- 插入项目经历数据
INSERT INTO t_rc_resume_project (id, resume_id, name, role, start_time, end_time, description, technology, achievement, created_at, updated_at) VALUES
(1, 1, '电商微服务重构项目', '核心开发', '2020-01-01', '2020-12-31', '将单体架构重构为微服务架构', 'Go,Docker,Kubernetes,Redis', '系统性能提升3倍，可维护性大幅提升', NOW(), NOW()),
(2, 2, '游戏实时对战系统', '技术负责人', '2019-03-01', '2020-02-28', '开发支持千万用户的实时对战系统', 'Java,Netty,Redis,MySQL', '支持同时在线用户数达到500万', NOW(), NOW()),
(3, 3, '搜索引擎优化项目', '架构师', '2017-06-01', '2018-12-31', '优化搜索引擎性能和准确率', 'Java,Elasticsearch,Kafka,HBase', '搜索响应时间降低40%，准确率提升15%', NOW(), NOW()),
(4, 4, '移动端游戏开发', '前端负责人', '2021-01-01', '2021-12-31', '开发H5移动端游戏', 'Vue.js,Canvas,WebGL,Node.js', '游戏上线3个月用户突破100万', NOW(), NOW()),
(5, 5, '云原生平台建设', '技术专家', '2018-01-01', '2019-12-31', '构建企业级云原生平台', 'Kubernetes,Docker,Istio,Prometheus', '平台支撑公司所有业务，节省运维成本60%', NOW(), NOW()),
(6, 6, '智能推荐系统', '产品负责人', '2020-01-01', '2021-06-30', '设计并实现个性化推荐系统', '机器学习,协同过滤,用户画像', '推荐点击率提升25%，用户活跃度提升30%', NOW(), NOW()),
(7, 7, 'CI/CD平台建设', '项目经理', '2019-01-01', '2020-12-31', '建设公司级CI/CD平台', 'Jenkins,GitLab,Docker,Ansible', '开发效率提升50%，发布故障率降低80%', NOW(), NOW()),
(8, 8, '移动应用UI重设计', '主设计师', '2021-01-01', '2021-08-31', '重新设计移动应用界面', 'Sketch,Figma,Principle,After Effects', '用户满意度从3.2提升到4.6', NOW(), NOW()),
(9, 9, '用户行为分析平台', '数据科学家', '2020-06-01', '2021-12-31', '构建用户行为分析和预测平台', 'Python,TensorFlow,Spark,Kafka', '预测准确率达到85%，为业务决策提供强力支撑', NOW(), NOW()),
(10, 10, '自动化测试平台', '测试工程师', '2022-01-01', '2022-12-31', '开发自动化测试平台', 'Selenium,Python,Pytest,Allure', '测试效率提升60%，bug发现率提升40%', NOW(), NOW());

-- 插入简历附件数据
INSERT INTO t_rc_resume_attachment (id, resume_id, file_name, file_url, file_size, file_type, status, created_at, updated_at) VALUES
(1, 1, '张三_简历.pdf', 'https://example.com/resumes/zhangsan_resume.pdf', 1024000, 'pdf', 1, NOW(), NOW()),
(2, 2, '李四_简历.docx', 'https://example.com/resumes/lisi_resume.docx', 856000, 'docx', 1, NOW(), NOW()),
(3, 3, '王五_简历.pdf', 'https://example.com/resumes/wangwu_resume.pdf', 1256000, 'pdf', 1, NOW(), NOW()),
(4, 4, '赵六_简历.pdf', 'https://example.com/resumes/zhaoliu_resume.pdf', 945000, 'pdf', 1, NOW(), NOW()),
(5, 5, '钱七_简历.pdf', 'https://example.com/resumes/qianqi_resume.pdf', 1345000, 'pdf', 1, NOW(), NOW()),
(6, 6, '孙八_简历.docx', 'https://example.com/resumes/sunba_resume.docx', 723000, 'docx', 1, NOW(), NOW()),
(7, 7, '周九_简历.pdf', 'https://example.com/resumes/zhoujiu_resume.pdf', 1123000, 'pdf', 1, NOW(), NOW()),
(8, 8, '吴十_简历.pdf', 'https://example.com/resumes/wushi_resume.pdf', 834000, 'pdf', 1, NOW(), NOW()),
(9, 9, '郑十一_简历.pdf', 'https://example.com/resumes/zhengshiyi_resume.pdf', 1567000, 'pdf', 1, NOW(), NOW()),
(10, 10, '冯十二_简历.docx', 'https://example.com/resumes/fengshier_resume.docx', 612000, 'docx', 1, NOW(), NOW());

-- 插入职位数据 (修正列名和数据格式)
INSERT INTO t_rc_job (id, name, company_id, job_skill, job_salary, job_salary_max, job_describe, job_location, job_expire_time, status, job_type, job_category, job_experience, job_education, job_benefit, job_contact, delete_status, job_source, create_time, update_time, view_count, apply_count, priority, tags, remote_type, remote_desc, remote_ratio, benefits, benefit_desc) VALUES
(1, '高级Go开发工程师', 2001, 'Go,Docker,Kubernetes,MySQL', 20000, 30000, '负责后端服务开发，要求有微服务经验', '北京', '2024-12-31 23:59:59', 2, 1, '技术', '3-5年', '本科', '五险一金,年终奖,股票期权', 'hr@example.com', 0, '官网', NOW(), NOW(), 245, 18, 1, '["急招","福利好"]', 2, '每周可远程3天', 60, '[1,2,7]', '额外商业保险'),
(2, 'Java架构师', 2002, 'Java,Spring,微服务,分布式', 35000, 50000, '负责系统架构设计，要求有大型项目经验', '上海', '2024-11-30 23:59:59', 2, 1, '技术', '5-8年', '本科', '五险一金,年终奖,期权激励', 'tech@company.com', 0, '官网', NOW(), NOW(), 312, 25, 2, '["高薪","发展好"]', 3, '全远程办公', 100, '[1,2,3,7]', '年度体检,健身补贴'),
(3, '前端开发工程师', 2003, 'Vue.js,React,TypeScript,Webpack', 15000, 25000, '负责前端界面开发，熟悉现代前端框架', '深圳', '2024-10-31 23:59:59', 2, 1, '技术', '1-3年', '本科', '五险一金,餐补,交通补助', 'frontend@tech.com', 0, '官网', NOW(), NOW(), 189, 32, 1, '["成长快","团队好"]', 2, '灵活办公时间', 40, '[1,5,6]', '免费午餐'),
(4, '产品经理', 2004, '产品设计,用户研究,数据分析,项目管理', 20000, 35000, '负责产品规划和设计，要求有B端产品经验', '杭州', '2024-09-30 23:59:59', 2, 1, '产品', '3-5年', '本科', '五险一金,年终奖,培训机会', 'pm@product.com', 0, '官网', NOW(), NOW(), 167, 14, 1, '["发展好","学习多"]', 1, '办公室办公', 0, '[1,2,4]', '丰富培训资源'),
(5, 'DevOps工程师', 2005, 'Docker,Kubernetes,Jenkins,AWS', 25000, 40000, '负责持续集成和部署，熟悉云原生技术', '广州', '2024-08-31 23:59:59', 2, 1, '技术', '3-5年', '本科', '五险一金,年终奖,技术津贴', 'devops@cloud.com', 0, '官网', NOW(), NOW(), 203, 19, 1, '["技术强","挑战大"]', 4, '弹性工作制', 80, '[1,2,4,8]', '技术书籍报销'),
(6, 'UI设计师', 2006, 'Sketch,Figma,Photoshop,UI设计', 12000, 20000, '负责产品界面设计，有移动端设计经验优先', '成都', '2024-07-31 23:59:59', 2, 1, '设计', '2-4年', '本科', '五险一金,设计津贴,创意奖金', 'design@creative.com', 0, '官网', NOW(), NOW(), 134, 21, 1, '["创意好","氛围佳"]', 2, '混合办公', 50, '[1,4,5]', '设计软件全报销'),
(7, '数据科学家', 2007, 'Python,机器学习,TensorFlow,Spark', 30000, 45000, '负责数据挖掘和机器学习模型开发', '西安', '2024-06-30 23:59:59', 2, 1, '技术', '3-6年', '硕士', '五险一金,科研津贴,论文奖励', 'data@ai.com', 0, '官网', NOW(), NOW(), 278, 12, 2, '["前沿技术","科研强"]', 3, '全远程', 100, '[1,2,4,9]', '学术会议支持'),
(8, '测试工程师', 2008, 'Selenium,自动化测试,Python,JMeter', 10000, 18000, '负责软件测试，熟悉自动化测试工具', '南京', '2024-05-31 23:59:59', 2, 1, '技术', '1-3年', '本科', '五险一金,项目奖金', 'qa@quality.com', 0, '官网', NOW(), NOW(), 156, 28, 1, '["学习多","稳定性好"]', 1, '办公室办公', 0, '[1,4]', '项目完成奖金'),
(9, 2009, '运营专员', 2009, '内容运营,数据分析,用户增长,活动策划', 8000, 15000, '负责产品运营工作，要求有用户增长经验', '武汉', '2024-04-30 23:59:59', 2, 1, '运营', '1-2年', '本科', '五险一金,绩效奖金', 'operation@growth.com', 0, '官网', NOW(), NOW(), 98, 35, 1, '["成长快","空间大"]', 2, '部分远程', 30, '[1,5]', '绩效奖金丰厚'),
(10, 2010, '销售经理', 2010, '销售技巧,客户管理,商务谈判,团队管理', 15000, 30000, '负责企业客户销售，有B端销售经验优先', '青岛', '2024-03-31 23:59:59', 2, 1, '销售', '3-5年', '本科', '五险一金,高额提成,出行补助', 'sales@business.com', 0, '官网', NOW(), NOW(), 87, 16, 1, '["高提成","发展好"]', 1, '需要出差', 0, '[1,6]', '销售提成无上限');

-- 插入职位申请数据
INSERT INTO t_rc_job_apply (id, job_id, company_id, user_id, resume_id, apply_time, apply_progress, reason, status, create_time, update_time) VALUES
(1, 1, 2001, 1001, 1, '2024-01-15 10:30:00', '面试中', '', 1, NOW(), NOW()),
(2, 2, 2002, 1003, 3, '2024-01-16 14:20:00', '已通过', '', 1, NOW(), NOW()),
(3, 3, 2003, 1004, 4, '2024-01-17 09:15:00', '待筛选', '', 1, NOW(), NOW()),
(4, 4, 2004, 1006, 6, '2024-01-18 16:45:00', '已拒绝', '经验不符合要求', 1, NOW(), NOW()),
(5, 5, 2005, 1007, 7, '2024-01-19 11:30:00', '面试通过', '', 1, NOW(), NOW()),
(6, 6, 2006, 1008, 8, '2024-01-20 13:20:00', '待面试', '', 1, NOW(), NOW()),
(7, 7, 2007, 1009, 9, '2024-01-21 15:10:00', 'Offer已发', '', 1, NOW(), NOW()),
(8, 8, 2008, 1010, 10, '2024-01-22 10:45:00', '简历筛选', '', 1, NOW(), NOW()),
(9, 1, 2001, 1002, 2, '2024-01-23 14:30:00', '已拒绝', '技术栈不匹配', 1, NOW(), NOW()),
(10, 9, 2009, 1005, 5, '2024-01-24 09:20:00', '待筛选', '', 1, NOW(), NOW());

-- 插入职位统计数据
INSERT INTO t_rc_job_statistics (id, job_id, company_id, view_count, apply_count, last_view_time, last_apply_time, conversion_rate, created_at, updated_at) VALUES
(1, 1, 2001, 245, 18, '2024-01-24 16:30:00', '2024-01-23 14:30:00', 7.35, NOW(), NOW()),
(2, 2, 2002, 312, 25, '2024-01-24 15:45:00', '2024-01-16 14:20:00', 8.01, NOW(), NOW()),
(3, 3, 2003, 189, 32, '2024-01-24 17:20:00', '2024-01-17 09:15:00', 16.93, NOW(), NOW()),
(4, 4, 2004, 167, 14, '2024-01-24 14:10:00', '2024-01-18 16:45:00', 8.38, NOW(), NOW()),
(5, 5, 2005, 203, 19, '2024-01-24 13:25:00', '2024-01-19 11:30:00', 9.36, NOW(), NOW()),
(6, 6, 2006, 134, 21, '2024-01-24 12:40:00', '2024-01-20 13:20:00', 15.67, NOW(), NOW()),
(7, 7, 2007, 278, 12, '2024-01-24 18:15:00', '2024-01-21 15:10:00', 4.32, NOW(), NOW()),
(8, 8, 2008, 156, 28, '2024-01-24 11:55:00', '2024-01-22 10:45:00', 17.95, NOW(), NOW()),
(9, 9, 2009, 98, 35, '2024-01-24 16:20:00', '2024-01-24 09:20:00', 35.71, NOW(), NOW()),
(10, 10, 2010, 87, 16, '2024-01-24 10:30:00', '2024-01-20 11:30:00', 18.39, NOW(), NOW());

-- 插入职位收藏数据 (修正列名，移除不存在的字段)
INSERT INTO t_rc_job_favorite (id, user_id, job_id, create_time, update_time) VALUES
(1, 1001, 2, '2024-01-15 10:30:00', NOW()),
(2, 1001, 7, '2024-01-16 14:20:00', NOW()),
(3, 1002, 1, '2024-01-17 09:15:00', NOW()),
(4, 1003, 5, '2024-01-18 16:45:00', NOW()),
(5, 1004, 3, '2024-01-19 11:30:00', NOW()),
(6, 1005, 9, '2024-01-20 13:20:00', NOW()),
(7, 1006, 4, '2024-01-21 15:10:00', NOW()),
(8, 1007, 5, '2024-01-22 10:45:00', NOW()),
(9, 1008, 6, '2024-01-23 14:30:00', NOW()),
(10, 1009, 7, '2024-01-24 09:20:00', NOW());

-- 插入简历交互数据
INSERT INTO t_rc_resume_interaction (id, resume_id, user_id, type, created_at, updated_at) VALUES
(1, 1, 2001, 'view', '2024-01-15 10:30:00', NOW()),
(2, 2, 2002, 'view', '2024-01-16 14:20:00', NOW()),
(3, 3, 2003, 'favorite', '2024-01-17 09:15:00', NOW()),
(4, 4, 2004, 'view', '2024-01-18 16:45:00', NOW()),
(5, 5, 2005, 'favorite', '2024-01-19 11:30:00', NOW()),
(6, 6, 2006, 'view', '2024-01-20 13:20:00', NOW()),
(7, 7, 2007, 'view', '2024-01-21 15:10:00', NOW()),
(8, 8, 2008, 'favorite', '2024-01-22 10:45:00', NOW()),
(9, 9, 2009, 'view', '2024-01-23 14:30:00', NOW()),
(10, 10, 2010, 'view', '2024-01-24 09:20:00', NOW());

-- 插入通知模板数据 (修正列名)
INSERT INTO t_rc_notification_template (id, code, title, content, type, user_types, channels, is_active, create_time, update_time) VALUES
(1, 'job_apply_success', '职位投递成功', '您已成功投递{{jobName}}职位，请耐心等待HR回复。', 1, '[1]', 1, true, NOW(), NOW()),
(2, 'job_apply_review', '简历已通过筛选', '恭喜！您投递的{{jobName}}职位简历已通过初筛，HR将尽快与您联系。', 2, '[1]', 1, true, NOW(), NOW()),
(3, 'job_apply_reject', '投递结果通知', '很遗憾，您投递的{{jobName}}职位暂未通过筛选，感谢您的关注。', 2, '[1]', 1, true, NOW(), NOW()),
(4, 'interview_invite', '面试邀请', '您好！{{companyName}}邀请您参加{{jobName}}职位的面试，时间：{{interviewTime}}。', 3, '[1]', 3, true, NOW(), NOW()),
(5, 'offer_send', 'Offer通知', '恭喜！{{companyName}}向您发送了{{jobName}}职位的Offer，请及时查看。', 2, '[1]', 1, true, NOW(), NOW()),
(6, 'resume_view', '简历被查看', '您的简历被{{companyName}}查看了，快去看看吧！', 4, '[1]', 1, true, NOW(), NOW()),
(7, 'job_expire_remind', '职位即将到期', '您发布的{{jobName}}职位将在3天后到期，请及时续期。', 4, '[2]', 1, true, NOW(), NOW()),
(8, 'system_maintain', '系统维护通知', '系统将于{{maintainTime}}进行维护，预计耗时{{duration}}。', 4, '[1,2,3,4]', 1, true, NOW(), NOW()),
(9, 'new_job_recommend', '为您推荐新职位', '根据您的简历，为您推荐{{jobName}}职位，快来看看吧！', 1, '[1]', 1, true, NOW(), NOW()),
(10, 'profile_complete', '完善简历信息', '完善您的简历信息，获得更多面试机会！', 4, '[1]', 1, true, NOW(), NOW());

-- 插入通知数据 (修正列名)
INSERT INTO t_rc_notification (id, user_id, user_type, type, title, content, channels, template_id, variables, is_read, create_time, update_time) VALUES
(1, 1001, 1, 1, '投递成功', '您已成功投递Go开发工程师职位，请耐心等待HR回复。', 1, 'job_apply_success', '{"jobName": "Go开发工程师"}', false, '2024-01-15 10:31:00', NOW()),
(2, 1003, 1, 2, '简历通过筛选', '恭喜！您投递的Java架构师职位简历已通过初筛。', 1, 'job_apply_review', '{"jobName": "Java架构师"}', false, '2024-01-16 14:21:00', NOW()),
(3, 1007, 1, 3, '面试邀请', 'DevOps工程师职位面试邀请，请准备面试。', 3, 'interview_invite', '{"jobName": "DevOps工程师", "companyName": "某公司"}', true, '2024-01-19 11:31:00', NOW()),
(4, 1009, 1, 2, 'Offer通知', '恭喜！数据科学家职位Offer已发送，请及时查看。', 1, 'offer_send', '{"jobName": "数据科学家"}', false, '2024-01-21 15:11:00', NOW()),
(5, 2001, 2, 1, '新的简历投递', '有候选人投递了您发布的Go开发工程师职位。', 1, '', '{"jobName": "Go开发工程师"}', true, '2024-01-15 10:31:00', NOW()),
(6, 2006, 2, 4, '简历被收藏', '您发布的UI设计师职位被用户收藏。', 1, '', '{"jobName": "UI设计师"}', false, '2024-01-23 14:31:00', NOW()),
(7, 1004, 1, 4, '职位推荐', '根据您的简历，为您推荐前端开发工程师职位。', 1, 'new_job_recommend', '{"jobName": "前端开发工程师"}', false, '2024-01-20 09:00:00', NOW()),
(8, 1010, 1, 4, '简历完善提醒', '完善您的简历信息，获得更多面试机会！', 1, 'profile_complete', '{"completeness": "60%"}', false, '2024-01-22 10:00:00', NOW()),
(9, 2007, 2, 4, '职位即将到期', '您发布的数据科学家职位将在3天后到期。', 1, 'job_expire_remind', '{"jobName": "数据科学家"}', false, '2024-01-24 09:00:00', NOW()),
(10, 1005, 1, 4, '系统维护通知', '系统将于明天凌晨2点进行维护，预计2小时。', 1, 'system_maintain', '{"maintainTime": "明天凌晨2点", "duration": "2小时"}', true, '2024-01-24 18:00:00', NOW());

COMMIT;
