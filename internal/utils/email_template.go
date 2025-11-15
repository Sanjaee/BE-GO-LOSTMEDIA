package utils

import (
	"fmt"
	"time"
)

// GetVerificationEmailTemplate returns HTML template for verification email with OTP
func GetVerificationEmailTemplate(otpCode string) string {
	year := time.Now().Year()
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Verify Your Email - LostMedia</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px;">
    <div style="background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); padding: 30px; text-align: center; border-radius: 10px 10px 0 0;">
        <h1 style="color: #ffffff; margin: 0; font-size: 28px;">LostMedia</h1>
    </div>
    
    <div style="background-color: #ffffff; padding: 40px; border: 1px solid #e0e0e0; border-top: none; border-radius: 0 0 10px 10px;">
        <h2 style="color: #333; margin-top: 0;">Verify Your Email Address</h2>
        
        <p style="color: #666; font-size: 16px;">
            Hello,
        </p>
        
        <p style="color: #666; font-size: 16px;">
            Thank you for registering with <strong>LostMedia</strong>! We're excited to have you on board.
        </p>
        
        <p style="color: #666; font-size: 16px;">
            To complete your registration and start using your account, please verify your email address using the OTP code below:
        </p>
        
        <div style="text-align: center; margin: 40px 0;">
            <div style="background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); color: #ffffff; padding: 20px; border-radius: 10px; display: inline-block;">
                <p style="margin: 0; font-size: 12px; letter-spacing: 2px; text-transform: uppercase; opacity: 0.9;">Your Verification Code</p>
                <p style="margin: 10px 0 0 0; font-size: 36px; font-weight: bold; letter-spacing: 8px; font-family: 'Courier New', monospace;">
                    %s
                </p>
            </div>
        </div>
        
        <p style="color: #666; font-size: 14px; margin-top: 30px; text-align: center;">
            Enter this code on the verification page to complete your registration.
        </p>
        
        <div style="margin-top: 40px; padding-top: 20px; border-top: 1px solid #e0e0e0;">
            <p style="color: #999; font-size: 12px; margin: 5px 0;">
                <strong>Important:</strong> This OTP code will expire in <strong>10 minutes</strong>.
            </p>
            <p style="color: #999; font-size: 12px; margin: 5px 0;">
                If you didn't create an account with LostMedia, please ignore this email.
            </p>
        </div>
    </div>
    
    <div style="text-align: center; margin-top: 20px; color: #999; font-size: 12px;">
        <p>© %d LostMedia. All rights reserved.</p>
    </div>
</body>
</html>`, otpCode, year)
}

// GetPasswordResetEmailTemplate returns HTML template for password reset email
func GetPasswordResetEmailTemplate(resetURL string) string {
	year := time.Now().Year()
	return fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Reset Your Password - LostMedia</title>
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px;">
    <div style="background: linear-gradient(135deg, #f093fb 0%%, #f5576c 100%%); padding: 30px; text-align: center; border-radius: 10px 10px 0 0;">
        <h1 style="color: #ffffff; margin: 0; font-size: 28px;">LostMedia</h1>
    </div>
    
    <div style="background-color: #ffffff; padding: 40px; border: 1px solid #e0e0e0; border-top: none; border-radius: 0 0 10px 10px;">
        <h2 style="color: #333; margin-top: 0;">Reset Your Password</h2>
        
        <p style="color: #666; font-size: 16px;">
            Hello,
        </p>
        
        <p style="color: #666; font-size: 16px;">
            We received a request to reset your password for your <strong>LostMedia</strong> account.
        </p>
        
        <p style="color: #666; font-size: 16px;">
            Click the button below to reset your password:
        </p>
        
        <div style="text-align: center; margin: 40px 0;">
            <a href="%s" style="background: linear-gradient(135deg, #f093fb 0%%, #f5576c 100%%); color: #ffffff; padding: 15px 40px; text-decoration: none; border-radius: 5px; display: inline-block; font-weight: bold; font-size: 16px;">
                Reset Password
            </a>
        </div>
        
        <p style="color: #666; font-size: 14px; margin-top: 30px;">
            Or copy and paste this link into your browser:
        </p>
        
        <p style="color: #f5576c; font-size: 12px; word-break: break-all; background-color: #f5f5f5; padding: 10px; border-radius: 5px;">
            %s
        </p>
        
        <div style="margin-top: 40px; padding-top: 20px; border-top: 1px solid #e0e0e0;">
            <p style="color: #999; font-size: 12px; margin: 5px 0;">
                <strong>Important:</strong> This password reset link will expire in <strong>1 hour</strong>.
            </p>
            <p style="color: #999; font-size: 12px; margin: 5px 0;">
                If you didn't request a password reset, please ignore this email. Your password will remain unchanged.
            </p>
            <p style="color: #999; font-size: 12px; margin: 5px 0;">
                For security reasons, please don't share this link with anyone.
            </p>
        </div>
    </div>
    
    <div style="text-align: center; margin-top: 20px; color: #999; font-size: 12px;">
        <p>© %d LostMedia. All rights reserved.</p>
    </div>
</body>
</html>`, resetURL, resetURL, year)
}
