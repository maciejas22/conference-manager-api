INSERT INTO public.terms_of_service (introduction, acknowledgement)
VALUES 
('Welcome to our service. Please read these terms carefully.', 'By using this service, you agree to our terms and conditions.');

INSERT INTO public.sections (terms_of_service_id, title, content)
VALUES 
(1, 'Introduction', 'This section introduces the terms of service.'),
(1, 'Eligibility', 'This section describes who is eligible to use the service.'),
(1, 'User Responsibilities', 'This section outlines the responsibilities of users.'),
(1, 'Privacy Policy', 'This section explains how user data is handled.'),
(1, 'Payment Terms', 'This section details the payment terms for using the service.'),
(1, 'Termination', 'This section describes the conditions for terminating the service.'),
(1, 'Dispute Resolution', 'This section explains how disputes are resolved.'),
(1, 'Changes to Terms', 'This section outlines how changes to the terms are made.');

INSERT INTO public.subsections (section_id, title, content)
VALUES 
(1, 'Purpose', 'This subsection explains the purpose of the introduction.'),
(2, 'Age Requirement', 'Users must be at least 18 years old.'),
(3, 'Account Security', 'Users must keep their account information secure.'),
(4, 'Data Collection', 'We collect data to improve our services.'),
(5, 'Refund Policy', 'Refunds are available under certain conditions.'),
(6, 'Account Termination', 'Accounts may be terminated for violations.'),
(7, 'Arbitration Process', 'Disputes will be resolved through arbitration.'),
(8, 'Notification of Changes', 'Users will be notified of changes to the terms.');
