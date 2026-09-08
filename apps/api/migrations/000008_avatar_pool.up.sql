-- The avatar fallback pool.
--
-- 94.3% of NextMoe accounts (119,242 of 126,383 on 2026-09-08) have no avatar
-- image, so whatever renders in their place IS the default avatar of the whole
-- ecosystem -- tens of millions of image requests a day. Until now every site
-- computed that image itself from a hardcoded `sticker.kungal.com/stickers/
-- KUNgal{set}/{n}.webp` URL, a path that encoded a POSITION IN A MUTABLE
-- COLLECTION and died the day this site stopped serving static files.
--
-- The replacement is a curated, slotted pool served from here as ready
-- content-addressed CDN URLs (see the avatar-pool endpoint). Three properties
-- matter and each is enforced below:
--
--  1. `slot` is a FIXED-LENGTH ARRAY INDEX, not a rank. Consumers pick with
--     `hash(name) % len(pool)`, so appending or removing an entry reshuffles
--     every user's identity. Replacing the sticker in one slot moves only the
--     users who hash to that slot. Edit slots; do not resize the array.
--  2. Membership is `slot IS NOT NULL`, and the unique index makes two
--     stickers in one slot impossible.
--  3. Only official, published packs may contribute -- enforced in the query,
--     not here, so that a hidden pack drops out of the pool the moment it is
--     hidden. The default avatar of 119k people must never be settable by an
--     arbitrary uploader.
--
-- The 64 seeded here were chosen by eye from all 498 official stickers,
-- rendered as 32px circles. The criterion is NOT "newest": it is "reads as a
-- face inside a circle" -- centred subject, no burned-in text, enough contrast
-- at 32px -- plus spread across characters and hair colours, because long runs
-- of the same character are the one real flaw in the source material.

ALTER TABLE sticker ADD COLUMN avatar_pool_slot SMALLINT;

CREATE UNIQUE INDEX sticker_avatar_pool_slot_idx
    ON sticker (avatar_pool_slot) WHERE avatar_pool_slot IS NOT NULL;

-- Seeded by image_hash, not by (pack, position): the hash IS the identity of
-- the bytes, and it is the one key that survives a repack or a reorder.
UPDATE sticker AS s
SET avatar_pool_slot = v.slot
FROM (VALUES
    (0, 'c697c3b779ad6c034a538516395a678f068afc1b24203e57d516f71bb09f070f'),  -- pack 1 #2
    (1, '2913c14bb5ffcde877a936520e2c1585a84dbba0afaf2d388058c22203d62b6c'),  -- pack 1 #10
    (2, 'dc8fd453f42d3edd7d51ac8a4ae1c7b698fb24d3166628d3560c5d4386bbf320'),  -- pack 1 #24
    (3, 'dcf153b812c8b62c8b93a75077dd38e2272745217dc7a1629c66a715e15660e5'),  -- pack 1 #25
    (4, 'a538294053a3a87ec08b21d07871d83b1e633c597cfd20188e954b205da457b1'),  -- pack 1 #27
    (5, 'd601f9745db81e72aeb0e153052ef1faf21427e90d956c2be55ab093fa817cee'),  -- pack 1 #40
    (6, '90d3ab0d3b04f5ff814d134c784cffb7b3942484ade8de19c8a11c68246f1c68'),  -- pack 1 #43
    (7, '2c6787e02437965dcc0588918b9fd9fe8246f8ab359259f2754743af0e66f716'),  -- pack 1 #46
    (8, '6016317952f8a945543208666069388cde91014ce07675050adb2557c542838c'),  -- pack 1 #51
    (9, '5ea8108bcab620e35ddb9dea37c19e9f1483cd2613b575d347b4e1afe040b9db'),  -- pack 1 #71
    (10, '540b51c8ebab0b47a24d47f2373b9e1a38bc07d2a8bb91ad2e46f632d30fa7a4'),  -- pack 1 #77
    (11, 'bf0ac1af8f1a1bc8983b187b66a93c1ccac43a35900d490ae01946fabfa3948d'),  -- pack 2 #3
    (12, '0da696d45c54aa2c8dcac3b26ab83c6774b655af1ab45faba5aec160b008f774'),  -- pack 2 #6
    (13, '549a84f2fa2b1d9c92b90183f0d33930292fd9b705c6ceb109f9f7cca61b139d'),  -- pack 2 #10
    (14, '0e63093d0255a90960f8fd29d80c80c35647e4ddc3e5c0ee433a760e0e890904'),  -- pack 2 #16
    (15, 'e26a73b9c36e9f19a25da7525ced393b1091bf81fa02fcc9c0ff4911b21f6810'),  -- pack 2 #29
    (16, '94b3ecb8a46f65584fea5e81d701c99bc498b5c1ee496e68ab7d97ceb35cb44f'),  -- pack 2 #35
    (17, 'b190356e1faf2a2a26e6e39acae4a4a411e88178ecc4511245ed9380963eb6d6'),  -- pack 2 #41
    (18, '09706ce22ba70fca249740aecb6fc0e1cdc95554688b70d26fed9d63192e074e'),  -- pack 2 #48
    (19, 'f924c4fcef2d1c68aaa255a83f1a6c3129256a5e13bc137df6197abcc10f7422'),  -- pack 2 #56
    (20, 'f1d8a8e734bedc4810d31a1df9b752df6b0f857cd93065c6411b22cd956c3613'),  -- pack 2 #61
    (21, '78356e14f95a941e27fcad3c338940f2911e987c28347625a14b50e9c457436a'),  -- pack 2 #66
    (22, 'b8a889468a1c6c0062157e0ad5c19a47c80d14cd23ad911cf6ccffa9f67d4fc2'),  -- pack 3 #6
    (23, '377d3e657cd8304876ff2a019e4094a6ed845e0d65bd24c88ab038a69e6b5316'),  -- pack 3 #14
    (24, '9681d8bde368a520468bf8708f56b5738ceeab8762b0e05dc822f38b228efd11'),  -- pack 3 #22
    (25, '86b9fe00a8564ed36f89575d5740c1d7c5dd4a174d8807499c479f7c705d6132'),  -- pack 3 #29
    (26, '7e22c431d10ddac463ba32960c178eaec9268a267ee5e5c61e3dd4972476fa50'),  -- pack 3 #32
    (27, 'd5ebfa226aeaef306a1ebf6a17cf4d4e8b16571df17839d74fe66f62edcfad98'),  -- pack 3 #35
    (28, 'caf455d1dcaacfd2819f4a41e91b14778d92c2ca952001ada084de6cb03412e4'),  -- pack 3 #45
    (29, '1a05a35370fc2264a1ae53c1feca16278ab932585052e516630e2411abe40399'),  -- pack 3 #49
    (30, '72d89db6890d79908a78632a1405ade6497063d10d1c67c5a8d3af70d8bdb0ee'),  -- pack 3 #72
    (31, 'b2a2d9bbd7517b70e0a148c73a5499c3c860588f816b65949d587594bebbee60'),  -- pack 4 #1
    (32, 'e4f6008a1fc9f970b06e33d91a2bf1613737d52711a31afae281a5503d028af0'),  -- pack 4 #5
    (33, 'b7b70f9b4f5332dce6fee9447fb2d7e303795b6f4ea3ea8d31c5cbc40e934d73'),  -- pack 4 #13
    (34, 'a1732838d54ecf4258d495e93af6c7c88118a624d7073124e32ca5e7bfbdf53a'),  -- pack 4 #20
    (35, '82ba7425100eccec3dbc6347ca74efe200d1bf80170c37031f7a43279754e1a6'),  -- pack 4 #22
    (36, 'ed9465a6b545e16eef6df62fdc8ecfee4295e4ec0d9bae3214294456ac6b92df'),  -- pack 4 #38
    (37, '771b0d99e97ca392588ff09bba231c60fbed847ebcdd75fc2d5f339b5ddc42d2'),  -- pack 4 #41
    (38, '365100c2907c038269c4fb51fa54bc77f07a157bc93f38d428c2f4b057896487'),  -- pack 4 #59
    (39, '72ee4cd2b1bb26054ae6f9e62f42f012ea26a26d2352d5ee03a478bccd643c81'),  -- pack 4 #61
    (40, 'e93e3198bc55d47931150492660a7ac1584eb110c127e3d0306dee55f36ab99a'),  -- pack 4 #68
    (41, '5f4003c2ca61252e8a04a64015579041d8ae713d774dd8461fdb8712781c1703'),  -- pack 4 #79
    (42, '29a4fc319249ef3d40c3036d3e386fcf2308e4f85bc33546fe52c0da344848e7'),  -- pack 5 #17
    (43, '4963541e788beca8acaffc115d5f8e3ad4bfd1404dde4ca986fbc9c76dc912a2'),  -- pack 5 #22
    (44, 'ffd7f32c2d43848047d7373a66187c31b05ba64a396b6adc17a9baab9ada49f2'),  -- pack 5 #25
    (45, '83af6a3fe04e534cf50b984633d8c0087ed03f6f927148377f79a3076ff7e3d2'),  -- pack 5 #30
    (46, '298ff6242a4f80e4290363660a1375d1d4949c3961e6427ef60b37377b2a0cff'),  -- pack 5 #35
    (47, '18d3223bb64bc1a2a38b60c519bb666dc0d222b3ee2a48f115d4089060186cc9'),  -- pack 5 #37
    (48, '72e4f2adfdb5301325a51a89fdfc59b8c0adfc93f6ade4939f329bcb6edcdd88'),  -- pack 5 #39
    (49, 'ec6a1a324aced75ea5a28f4c3cfc813db8c06289e77e0078c3251e7a9cdaa052'),  -- pack 5 #64
    (50, '2f160616d59a51b294c578191f657ec49ea5ba78e6cee7d63387a52c4a54bba4'),  -- pack 5 #68
    (51, '915caccac223dbbab3829d593fdde0b07cab6af776ffe51b32ff4560c3cd19db'),  -- pack 6 #20
    (52, '4b3f944f4b262dc8dd152733db791f7ba4f4197c1587c17d6cecdde21e984389'),  -- pack 6 #21
    (53, 'c44c91cb036a7ee197acaa27cce5d9b5eb5d42c40b887e2e113a371f8ca8f974'),  -- pack 6 #25
    (54, '97cf96ba7c9e05797dbdcd012e257e4f659fb74527f631b91d57ffe2be7f21e7'),  -- pack 6 #41
    (55, 'f4ff43261e1d4754d797410d3657fbdcc586c7a5294d6f88485fe6daf302969b'),  -- pack 6 #47
    (56, 'd77d2e1edd96760faa4319ae727e0a77532ff4e2367b73ab6b853d2ee1788d39'),  -- pack 6 #51
    (57, 'ea71f5394ec86b525f00b1289048fca98a05b0396a4502fcf4965ad581359792'),  -- pack 6 #57
    (58, '6620f4d2cf32ae46b84481d5dd8c42fd85eb89394b494166c96e1d8395ab9657'),  -- pack 6 #59
    (59, '128deb5e43fbfd1aaf1b47f18e9a6b7497c4620cba90129b8cd1094348c36120'),  -- pack 6 #61
    (60, '8e415e9945d6c5113dd10bcc6687f91e38c883c1d4479552ae5b86fc82e54d02'),  -- pack 6 #65
    (61, '163a237ae4590e2c8c62d0df0144e3e2c1df1704819e0bc2b06449300a0a1235'),  -- pack 6 #71
    (62, 'b7f255c80c0aaecaff00eff4e3155dca05335d1ab71f6b7c5c21798273eef1da'),  -- pack 7 #2
    (63, '3e609eaf0c40913f5e1273acbbad7833740d105f31e1a8eb5c546c2426bcd10c')   -- pack 7 #15
) AS v(slot, image_hash)
WHERE s.image_hash = v.image_hash;
